package bots

import (
	"encoding/json"

	"fmt"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

type Version struct {
	w *core.Wrangler
}

func (b Version) Name() string {
	return "Versionista"
}

func (b Version) Description() string {
	return "Can look up a version across environments"
}

func (b Version) Help() hipchat.NotificationRequest {
	help := `
	%s - %s
	<br>- /bot version $app - Returns information about $app on all environments
	`
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       fmt.Sprintf(help, b.Name(), b.Description()),
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b Version) HandleMessage(req *webhook.Request) hipchat.NotificationRequest {
	app := req.Message()
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorGray,
		Message:       fmt.Sprintf("Versions for %s: [tst] 1025.255 [acc] 3636 [xpr] 335 [pro] 3636", app),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b Version) HandleConfig(wrangler *core.Wrangler, data json.RawMessage) error {
	b.w = wrangler
	return nil
}
