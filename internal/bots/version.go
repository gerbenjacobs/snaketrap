package bots

import (
	"encoding/json"

	"fmt"

	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
)

type Version struct{}

func (b Version) Name() string {
	return "Versionista"
}

func (b Version) Description() string {
	return "Can look up a version across environments"
}

func (b Version) Help() *hipchat.Response {
	help := `
	%s - %s
	<br>- /bot version $app - Returns information about $app on all environments
	`
	return hipchat.NewHelp(fmt.Sprintf(help, b.Name(), b.Description()))
}

func (b Version) HandleMessage(c *hipchat.Client, req *hipchat.Request) *hipchat.Response {
	app := req.GetWord(2)
	return hipchat.NewResponse(hipchat.COLOR_GRAY, fmt.Sprintf("Versions for %s: [tst] 1025.255 [acc] 3636 [xpr] 335 [pro] 3636", app))
}

func (b Version) HandleConfig(data json.RawMessage) error {
	return nil
}
