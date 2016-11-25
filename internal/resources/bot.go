package resources

import (
	"net/http"

	"encoding/json"
	"fmt"

	"bytes"

	"strings"

	"net/url"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-selfdiagnose"
	"github.com/gerbenjacobs/snaketrap/internal/bots"
	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

type BotConfig struct {
	Enabled bool            `json:"enabled"`
	Data    json.RawMessage `json:"data"`
}

var BotLookup = map[string]core.Bot{
	"sheriff": &bots.Sheriff{},
	"version": &bots.Version{},
}

type BotResource struct {
	wrangler *core.Wrangler
	bots     map[string]core.Bot
}

// Selfdiagnose
func (b BotResource) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	var w bytes.Buffer
	for botKey, bot := range b.bots {
		fmt.Fprintf(&w, "[%v] %v - %v<br>", botKey, bot.Name(), bot.Description())
	}

	result.Passed = true
	result.Reason = w.String()
}
func (b BotResource) Comment() string {
	return "Bots"
}

func NewBotResource(wrangler *core.Wrangler) (*BotResource, error) {
	// create Hipchat clients
	c := hipchat.NewClient(wrangler.ScopeAuth)
	b := hipchat.NewClient(wrangler.BotAuth)
	pUrl, err := url.Parse(wrangler.Url)
	if err != nil {
		return nil, err
	}
	c.BaseURL = pUrl
	b.BaseURL = pUrl
	wrangler.Client = c
	wrangler.SetBotClient(b)

	// Configure and return
	botMap, err := CreateAndConfigureBots(wrangler)
	if err != nil {
		return nil, err
	}
	return &BotResource{
		wrangler: wrangler,
		bots:     botMap,
	}, nil
}

func CreateAndConfigureBots(wrangler *core.Wrangler) (map[string]core.Bot, error) {
	var botConfig map[string]BotConfig
	if err := json.Unmarshal(wrangler.RawData["bots"], &botConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bot config: %v", err)
	}

	botMap := map[string]core.Bot{}
	for botName, bc := range botConfig {
		if bc.Enabled {
			// create
			if _, ok := BotLookup[botName]; !ok {
				return nil, fmt.Errorf("failed to find bot [%v]", botName)
			}
			myBot := BotLookup[botName]

			// configure
			if err := myBot.HandleConfig(wrangler, botConfig[botName].Data); err != nil {
				return nil, fmt.Errorf("failed to configure bot [%v]: %v", botName, err)
			}

			botMap[botName] = myBot
		}
	}

	return botMap, nil
}

// bind this resource to the container
func (b BotResource) Bind(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/bot").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(b.handleRequest))

	container.Add(ws)
}
func (b *BotResource) handleRequest(request *restful.Request, response *restful.Response) {
	// process request
	req := new(webhook.Request)
	err := request.ReadEntity(req)
	if err != nil {
		log15.Error("failed to read entity", "err", err)
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// pick bot and create notification
	botName := req.Bot()
	if "--help" == botName {
		// send general help
		response.WriteEntity(b.HelpMsg())
		return
	}

	var notification hipchat.NotificationRequest
	var use bool
	if bot, ok := b.bots[botName]; ok {
		if strings.Contains(req.Message(), "--help") {
			notification = bot.Help()
		} else {
			notification, use = bot.HandleMessage(req)
		}
		notification.From = bot.Name()
	} else {
		log15.Error("failed to handle request", "bot", botName, "msg", req.Message())
		notification = b.FailedMsg()
	}

	// reply :)
	log15.Debug("handled bot request", "from", req.Username(), "message", req.Message(), "replying", use, "notification", notification)
	if use {
		// send back notification
		response.WriteEntity(notification)
	} else {
		// received call, but nothing to reply
		response.WriteHeader(201)
	}
}

func (b BotResource) HelpMsg() hipchat.NotificationRequest {
	bs := []string{}
	for i := range b.bots {
		bs = append(bs, i)
	}

	help := fmt.Sprintf("Current active bots: %s<br>Use <strong>/bot $name --help</strong> for information per bot", strings.Join(bs, ", "))

	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       help,
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b BotResource) FailedMsg() hipchat.NotificationRequest {
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorRed,
		Message:       "I'm sorry, I don't understand that command.",
		Notify:        false,
		MessageFormat: "text",
	}
}
