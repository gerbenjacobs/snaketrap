package main

import (
	"net/http"

	"flag"

	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
	"github.com/gerbenjacobs/snaketrap/internal/resources"
	"github.com/inconshreveable/log15"
)

var (
	port   int
	hcUrl  string
	hcAuth string
	hcRoom int
)

func main() {
	// [*] Configuration
	flag.IntVar(&port, "port", 8080, "The port that the HTTP listener runs on")
	flag.StringVar(&hcUrl, "hipchat-url", "", "HipChat API endpoint, including version")
	flag.StringVar(&hcAuth, "hipchat-auth", "", "HipChat auth token")
	flag.IntVar(&hcRoom, "hipchat-room-id", 0, "HipChat room ID")
	flag.Parse()

	// [*] Create the container
	container := restful.DefaultContainer
	container.Router(restful.CurlyRouter{})

	// [*] Middleware; logging and instrumentation
	container.Filter(core.AccessLogger())

	// [*] Create the main app
	app := new(restful.WebService)
	app.Route(app.GET("/").To(rootRedirect))

	// [*] Self diagnose
	core.SetupSelfdiagnose(app)

	// Create *your* resources and bind them
	hc := hipchat.NewClient(hcUrl, hcAuth, hcRoom)
	bot := resources.NewBotResource(hc)
	bot.Bind(container)

	// [*] Enable swagger, after resources have been bound
	core.EnableSwagger(container)

	// [*] Bind the main app to the container
	container.Add(app)

	// [*] Run server
	log15.Info("Started listening on 0.0.0.0:" + strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), container)
	if err != nil {
		log15.Info("failed to listen and serve", "err", err)
	}
}

func rootRedirect(req *restful.Request, resp *restful.Response) {
	http.Redirect(resp.ResponseWriter, req.Request, "/internal/selfdiagnose.html", http.StatusPermanentRedirect)
}
