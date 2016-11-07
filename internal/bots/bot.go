package bots

import (
	"encoding/json"

	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
)

type BotConfig struct {
	Enabled bool            `json:"enabled"`
	Data    json.RawMessage `json:"data"`
}

type Bot interface {
	Name() string
	Description() string
	Help() *hipchat.Response
	HandleMessage(*hipchat.Client, *hipchat.Request) *hipchat.Response
	HandleConfig(json.RawMessage) error
}
