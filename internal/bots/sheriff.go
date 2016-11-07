package bots

import (
	"time"

	"encoding/json"

	"fmt"

	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
)

type SheriffConfig struct {
	users map[string]string `json:"users"`
}

type Sheriff struct {
	config SheriffConfig
}

func (b Sheriff) Name() string {
	return "Sheriff"
}

func (b Sheriff) Description() string {
	return "A bot that rotates users daily for the engineer on duty role a.k.a. sheriff"
}

func (b Sheriff) Help() *hipchat.Response {
	help := `
	%s - %s
	<br>- /bot sheriff <strong>next</strong> - Switches to next sheriff
	<br>- /bot sheriff <strong>previous</strong> - Switches to previous sheriff
	<br>- /bot sheriff <strong>away</strong> $user - Marks the $user as away
	<br>- /bot sheriff <strong>back</strong> $user - Marks the $user as back
	`
	return hipchat.NewHelp(fmt.Sprintf(help, b.Name(), b.Description()))
}

func (b Sheriff) HandleMessage(c *hipchat.Client, req *hipchat.Request) *hipchat.Response {
	go c.ChangeRoomTopic("Sheriff of the day - " + time.Now().String())
	return hipchat.NewResponse(hipchat.COLOR_YELLOW, "Sheriff of the day is: you!")
}

func (b Sheriff) HandleConfig(data json.RawMessage) error {
	var cfg SheriffConfig
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	b.config = cfg
	return nil
}
