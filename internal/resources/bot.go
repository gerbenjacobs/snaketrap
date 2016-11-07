package resources

import (
	"net/http"

	"encoding/json"
	"fmt"

	"bytes"

	"strings"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-selfdiagnose"
	"github.com/gerbenjacobs/snaketrap/internal/bots"
	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
	"github.com/inconshreveable/log15"
)

var BotLookup = map[string]bots.Bot{
	"sheriff": &bots.Sheriff{},
	"version": &bots.Version{},
}

type BotResource struct {
	client *hipchat.Client
	bots   map[string]bots.Bot
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

func NewBotResource(client *hipchat.Client, cfgMap map[string]json.RawMessage) (*BotResource, error) {
	botMap, err := CreateAndConfigureBots(cfgMap)
	if err != nil {
		return nil, err
	}
	return &BotResource{
		client: client,
		bots:   botMap,
	}, nil
}

func CreateAndConfigureBots(cfgMap map[string]json.RawMessage) (map[string]bots.Bot, error) {
	var botConfig map[string]bots.BotConfig
	if err := json.Unmarshal(cfgMap["bots"], &botConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bot config: %v", err)
	}

	botMap := map[string]bots.Bot{}
	for botName, bc := range botConfig {
		if bc.Enabled {
			// create
			if _, ok := BotLookup[botName]; !ok {
				return nil, fmt.Errorf("failed to find bot [%v]", botName)
			}
			myBot := BotLookup[botName]

			// configure
			if err := myBot.HandleConfig(botConfig[botName].Data); err != nil {
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
		Doc("Entrypoint for the bot").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(b.handleRequest).
		Doc("Handles incoming requests").
		Reads(hipchat.Request{}).
		Writes(hipchat.Response{}))

	container.Add(ws)
}
func (b *BotResource) handleRequest(request *restful.Request, response *restful.Response) {
	// process request
	hcReq := new(hipchat.Request)
	err := request.ReadEntity(hcReq)
	if err != nil {
		log15.Error("failed to read entity", "err", err)
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// pick bot and create notification
	botName := hcReq.Bot()
	if "--help" == botName {
		// send general help
		response.WriteEntity(b.HelpMsg())
		return
	}

	var notification *hipchat.Response
	if bot, ok := b.bots[botName]; ok {
		if "--help" == hcReq.GetWord(2) {
			notification = bot.Help()
		} else {
			notification = bot.HandleMessage(b.client, hcReq)
		}
		notification.From = bot.Name()
		notification.AttachTo = hcReq.Item.Message.ID
	} else {
		notification = hipchat.NewResponse(hipchat.COLOR_RED, "I'm sorry, I don't understand that command.")
		log15.Error("failed to handle request", "bot", botName, "msg", hcReq.Message())
	}

	// reply :)
	response.WriteEntity(notification)
}

func (b BotResource) HelpMsg() *hipchat.Response {
	bs := []string{}
	for i := range b.bots {
		bs = append(bs, i)
	}

	help := fmt.Sprintf("Current active bots: %s<br>Use <strong>/bot $name --help</strong> for information per bot", strings.Join(bs, ", "))

	return hipchat.NewHelp(help)
}
