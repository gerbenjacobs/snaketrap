package bots

import (
	"time"

	"encoding/json"

	"fmt"

	"github.com/tbruyelle/hipchat-go/hipchat"

	"github.com/gerbenjacobs/snaketrap/internal/core"
	"github.com/inconshreveable/log15"
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

func (b Sheriff) Help() hipchat.NotificationRequest {
	help := `
	%s - %s
	<br>- /bot sheriff <strong>next</strong> - Switches to next sheriff
	<br>- /bot sheriff <strong>previous</strong> - Switches to previous sheriff
	<br>- /bot sheriff <strong>away</strong> $user - Marks the $user as away
	<br>- /bot sheriff <strong>back</strong> $user - Marks the $user as back
	`
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       fmt.Sprintf(help, b.Name(), b.Description()),
		Notify:        false,
		MessageFormat: "html",
	}
}

func (b Sheriff) HandleMessage(c *core.HipchatConfig, req *hipchat.RoomMessageRequest) hipchat.NotificationRequest {
	go func() {
		_, err := c.Client.Room.SetTopic(c.DefaultRoom, "Sheriff of the day - "+time.Now().String())
		if err != nil {
			log15.Error("failed to set topic", "err", err)
		}
	}()
	return hipchat.NotificationRequest{
		Color:         hipchat.ColorYellow,
		Message:       "Sheriff of the day is: you!",
		Notify:        false,
		MessageFormat: "text",
	}
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
