package resources

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/gerbenjacobs/snaketrap/internal/bots"
	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
	"gopkg.in/inconshreveable/log15.v2"
)

type BotResource struct {
	client *hipchat.Client
	bots   map[string]bots.Bot
}

func NewBotResource(client *hipchat.Client) BotResource {
	return BotResource{
		client: client,
		bots: map[string]bots.Bot{
			"sheriff": bots.Sheriff{},
		},
	}
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

	ws.Route(ws.POST("/notify").To(b.notify).
		Doc("Sends a notification").
		Param(ws.QueryParameter("message", "Message to be sent").DataType("string")))

	container.Add(ws)
}

func (b *BotResource) notify(request *restful.Request, response *restful.Response) {
	notification := hipchat.NewNotify(hipchat.COLOR_PURPLE, request.QueryParameter("message"))
	notification.From = "/notify"
	err := b.client.SendNotification(notification)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	response.Write([]byte("OK"))
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

	// pick bot and creat notification
	var notification *hipchat.Response
	botName := hcReq.Bot()
	if bot, ok := b.bots[botName]; ok {
		notification = bot.HandleMessage(b.client, hcReq)
		notification.From = bot.Name()
	} else {
		notification = hipchat.NewResponse(hipchat.COLOR_RED, "I'm sorry, I don't understand that command.")
	}

	// reply :)
	response.WriteEntity(notification)
}
