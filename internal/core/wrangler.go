package core

import (
	"encoding/json"

	"io/ioutil"

	"os"

	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

// Wrangler is the object that holds the configuration, HipChat client
// and has several general methods
type Wrangler struct {
	URL         string `json:"url"`
	BotAuth     string `json:"bot_auth"`
	ScopeAuth   string `json:"scope_auth"`
	DefaultRoom string `json:"room_id"`
	Client      *hipchat.Client
	botClient   *hipchat.Client
	RawData     map[string]json.RawMessage
}

// Bot is the interface for the small bots used by the wrangler
type Bot interface {
	Name() string
	Description() string
	Help() hipchat.NotificationRequest
	HandleMessage(*webhook.Request) (hipchat.NotificationRequest, bool)
	HandleConfig(*Wrangler, json.RawMessage) error
	CurrentState() []byte
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
	w.botClient.Room.Notification(w.DefaultRoom, n)
}

func (w *Wrangler) SetState(b Bot) error {
	return ioutil.WriteFile(hashBotStateFile(b.Name()), b.CurrentState(), 0640)
}

func (w *Wrangler) GetState(b Bot) ([]byte, error) {
	f := hashBotStateFile(b.Name())
	data, err := ioutil.ReadFile(f)
	if _, ok := err.(*os.PathError); ok {
		// no file, initialize
		return w.initializeState(b, f)
	}
	return data, err
}

func (w *Wrangler) initializeState(b Bot, f string) ([]byte, error) {
	if err := os.MkdirAll(StateFileFolder, os.ModePerm); err != nil {
		return nil, err
	}
	if err := w.SetState(b); err != nil {
		log15.Warn("failed to create initial state file", "err", err, "bot", b.Name(), "filename", f)
		return nil, err
	}

	// created, call myself again
	return w.GetState(b)
}
