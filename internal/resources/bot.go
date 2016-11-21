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

type Bot interface {
	Name() string
	Description() string
	Help() hipchat.NotificationRequest
	HandleMessage(*webhook.Request) hipchat.NotificationRequest
	HandleConfig(*core.HipchatConfig, json.RawMessage) error
}

var BotLookup = map[string]Bot{
	"sheriff": &bots.Sheriff{},
	"version": &bots.Version{},
}

type BotResource struct {
	config *core.HipchatConfig
	bots   map[string]Bot
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

func NewBotResource(cfg *core.HipchatConfig, cfgMap map[string]json.RawMessage) (*BotResource, error) {
	// create Hipchat client
	c := hipchat.NewClient(cfg.Auth)
	pUrl, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, err
	}
	c.BaseURL = pUrl
	cfg.Client = c

	// Configure and return
	botMap, err := CreateAndConfigureBots(cfg, cfgMap)
	if err != nil {
		return nil, err
	}
	return &BotResource{
		config: cfg,
		bots:   botMap,
	}, nil
}

func CreateAndConfigureBots(c *core.HipchatConfig, cfgMap map[string]json.RawMessage) (map[string]Bot, error) {
	var botConfig map[string]BotConfig
	if err := json.Unmarshal(cfgMap["bots"], &botConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bot config: %v", err)
	}

	botMap := map[string]Bot{}
	for botName, bc := range botConfig {
		if bc.Enabled {
			// create
			if _, ok := BotLookup[botName]; !ok {
				return nil, fmt.Errorf("failed to find bot [%v]", botName)
			}
			myBot := BotLookup[botName]

			// configure
			if err := myBot.HandleConfig(c, botConfig[botName].Data); err != nil {
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
	if bot, ok := b.bots[botName]; ok {
		if strings.Contains(req.Message(), "--help") {
			notification = bot.Help()
		} else {
			notification = bot.HandleMessage(req)
		}
		notification.From = bot.Name()
	} else {
		log15.Error("failed to handle request", "bot", botName, "msg", req.Message())
		notification = b.FailedMsg()
	}

	// reply :)
	log15.Debug("handled bot request", "from", req.Username(), "message", req.Message(), "reply", notification)
	response.WriteEntity(notification)
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
