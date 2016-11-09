package main

import (
	"net/http"

	"os"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-selfdiagnose"
	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/resources"
	"github.com/inconshreveable/log15"
)

func main() {
	// [*] Configuration
	var addr string
	var hcConfig core.HipchatConfig
	cfgMap, err := core.ReadConfig(&addr, &hcConfig)
	if err != nil {
		log15.Error("failed to read config", "err", err)
		os.Exit(1)
	}

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
	bot, err := resources.NewBotResource(&hcConfig, cfgMap)
	if err != nil {
		log15.Error("failed to create BotResource", "err", err)
		os.Exit(1)
	}
	bot.Bind(container)
	selfdiagnose.Register(bot)

	// [*] Bind the main app to the container
	container.Add(app)

	// [*] Run server
	log15.Info("Started listening on " + addr)
	err = http.ListenAndServe(addr, container)
	if err != nil {
		log15.Info("failed to listen and serve", "err", err)
	}
}

func rootRedirect(req *restful.Request, resp *restful.Response) {
	http.Redirect(resp.ResponseWriter, req.Request, "/internal/selfdiagnose.html", http.StatusTemporaryRedirect)
}
