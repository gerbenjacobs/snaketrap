package core

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/gerbenjacobs/snaketrap/internal/webhook"
	"github.com/inconshreveable/log15"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

// HipchatConfig is the struct used to marshall the "hipchat" part
// of the JSON configuration
type HipchatConfig struct {
	URL         string `json:"url"`
	BotAuth     string `json:"bot_auth"`
	ScopeAuth   string `json:"scope_auth"`
	DefaultRoom string `json:"room_id"`
}

// Wrangler is the object that holds the configuration, HipChat client
// and has several general methods
type Wrangler struct {
	defaultRoom string
	scopeClient *hipchat.Client
	botClient   *hipchat.Client
	rawData     map[string]json.RawMessage
}

// Bot is the interface for the small bots used by the wrangler
type Bot interface {
	Name() string
	Description() string
	Help() Reply
	HandleMessage(*webhook.Request) Reply
	HandleConfig(*Wrangler, json.RawMessage) error
	CurrentState() []byte
}

type Reply struct {
	Notification hipchat.NotificationRequest
	Replying bool
}

func NewReply(n hipchat.NotificationRequest) Reply {
	return Reply{Notification: n, Replying: true}
}

func NoOpReply() Reply {
	return Reply{Notification: hipchat.NotificationRequest{}, Replying: false}
}

// NewSnaketrap reads the config.json and returns a bootstrapped Wrangler and listen address.
func NewSnaketrap(fileName string) (*Wrangler, string, error) {
	addr := ""
	jsData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, addr, err
	}

	// get raw data
	cfgMap := map[string]json.RawMessage{}
	if err = json.Unmarshal(jsData, &cfgMap); err != nil {
		return nil, addr, err
	}

	// marshall the addr used in main
	if err = json.Unmarshal(cfgMap["addr"], &addr); err != nil {
		return nil, addr, err
	}

	// try to marshall the known properties
	hcc := &HipchatConfig{}
	if err = json.Unmarshal(cfgMap["hipchat"], hcc); err != nil {
		return nil, addr, err
	}

	// bootstrap the wrangler object
	wrangler := &Wrangler{}
	if err = wrangler.bootstrap(hcc, cfgMap); err != nil {
		return nil, "", err
	}

	return wrangler, addr, nil
}

func (w *Wrangler) bootstrap(hcc *HipchatConfig, cfgMap map[string]json.RawMessage) error {
	// create HipChat clients
	c := hipchat.NewClient(hcc.ScopeAuth)
	b := hipchat.NewClient(hcc.BotAuth)
	pURL, err := url.Parse(hcc.URL)
	if err != nil {
		return err
	}
	c.BaseURL = pURL
	b.BaseURL = pURL

	// set up Wrangler
	w.defaultRoom = hcc.DefaultRoom
	w.scopeClient = c
	w.botClient = b
	w.rawData = cfgMap

	return nil
}

func (w *Wrangler) GetBotConfig() json.RawMessage {
	return w.rawData["bots"]
}

func (w *Wrangler) SendNotification(b Bot, n *hipchat.NotificationRequest) {
	n.From = b.Name()
	_, err := w.botClient.Room.Notification(w.defaultRoom, n)
	if err != nil {
		log15.Error("failed to send noticiation", "err", err, "bot", b.Name())
	}
}

func (w *Wrangler) SetTopic(b Bot, topic string) {
	_, err := w.scopeClient.Room.SetTopic(w.defaultRoom, topic)
	if err != nil {
		log15.Error("failed to set topic", "err", err, "bot", b.Name())
	}
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
