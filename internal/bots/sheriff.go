package bots

import (
	"time"

	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
)

type Sheriff struct{}

func (s Sheriff) Name() string {
	return "Sheriff"
}

func (s Sheriff) HandleMessage(c *hipchat.Client, req *hipchat.Request) *hipchat.Response {
	c.ChangeRoomTopic("Sheriff of the day - " + time.Now().String())
	return hipchat.NewResponse(hipchat.COLOR_YELLOW, "Sheriff of the day is: you!")
}
