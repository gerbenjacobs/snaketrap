package core

import (
	"encoding/json"

	"io/ioutil"

	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

// Wrangler is the object that holds the configuration, HipChat client
// and has several general methods
type Wrangler struct {
	Url         string `json:"url"`
	BotAuth     string `json:"bot_auth"`
	ScopeAuth   string `json:"scope_auth"`
	DefaultRoom string `json:"room_id"`
	Client      *hipchat.Client
	botClient   *hipchat.Client
	RawData     map[string]json.RawMessage
}

// Bot is the interface uses for the small bots used by the wrangler
type Bot interface {
	Name() string
	Description() string
	Help() hipchat.NotificationRequest
	HandleMessage(*webhook.Request) (hipchat.NotificationRequest, bool)
	HandleConfig(*Wrangler, json.RawMessage) error
}

func ReadConfig(addr *string, wrangler *Wrangler) error {
	jsData, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	// get raw data
	cfgMap := map[string]json.RawMessage{}
	err = json.Unmarshal(jsData, &cfgMap)
	if err != nil {
		return err

	}

	// marshall the addr used in main
	err = json.Unmarshal(cfgMap["addr"], &addr)
	if err != nil {
		return err
	}

	// try to marshall the known properties
	err = json.Unmarshal(cfgMap["hipchat"], &wrangler)
	if err != nil {
		return err
	}

	// add the rest of the raw data to the wrangler object
	wrangler.RawData = cfgMap

	return nil
}

func (w *Wrangler) SetBotClient(c *hipchat.Client) {
	w.botClient = c
}

func (w *Wrangler) SendNotification(b Bot, n *hipchat.NotificationRequest) {
	n.From = b.Name()
	log15.Info("Trying to send general Bot notification", "n", n)
	w.botClient.Room.Notification(w.DefaultRoom, n)
}
