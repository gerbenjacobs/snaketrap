package core

import (
	"encoding/json"

	"io/ioutil"

	"github.com/tbruyelle/hipchat-go/hipchat"
	"gopkg.in/inconshreveable/log15.v2"
)

// Wrangler is the object that holds the configuration, HipChat client
// and has several general methods
type Wrangler struct {
	Url         string `json:"url"`
	Auth        string `json:"auth"`
	DefaultRoom string `json:"roomId"`
	Client      *hipchat.Client
	RawData     map[string]json.RawMessage
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

func (w *Wrangler) sendNotification(n *hipchat.NotificationRequest) {
	log15.Info("Trying to send general Bot notification", "n", n)
}
