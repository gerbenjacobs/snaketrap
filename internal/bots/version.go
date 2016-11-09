package bots

import (
	"encoding/json"

	"fmt"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

type Version struct{}

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

func (b Version) HandleMessage(c *core.HipchatConfig, req *hipchat.RoomMessageRequest) hipchat.NotificationRequest {
	app := req.Message
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorGray,
		Message:       fmt.Sprintf("Versions for %s: [tst] 1025.255 [acc] 3636 [xpr] 335 [pro] 3636", app),
		Notify:        false,
		MessageFormat: "text",
	}
}

func (b Version) HandleConfig(data json.RawMessage) error {
	return nil
}
