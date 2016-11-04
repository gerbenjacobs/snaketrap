package bots

import "github.com/gerbenjacobs/snaketrap/internal/hipchat"

type Bot interface {
	Name() string
	HandleMessage(*hipchat.Client, *hipchat.Request) *hipchat.Response
}
