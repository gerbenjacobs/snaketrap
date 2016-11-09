package core

import (
	"encoding/json"

	"io/ioutil"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

type HipchatConfig struct {
	Url         string `json:"url"`
	Auth        string `json:"auth"`
	DefaultRoom string `json:"roomId"`
	Client      *hipchat.Client
}

func ReadConfig(addr *string, hcConfig *HipchatConfig) (cfgMap map[string]json.RawMessage, err error) {
	jsData, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsData, &cfgMap)
	if err != nil {
		return nil, err

	}

	err = json.Unmarshal(cfgMap["addr"], &addr)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfgMap["hipchat"], &hcConfig)
	if err != nil {
		return nil, err
	}

	return cfgMap, nil
}
