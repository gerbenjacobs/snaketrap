package core

import (
	"encoding/json"

	"io/ioutil"

	"github.com/gerbenjacobs/snaketrap/internal/hipchat"
)

func ReadConfig(port *int, hcClient *hipchat.Client) (cfgMap map[string]json.RawMessage, err error) {
	jsData, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsData, &cfgMap)
	if err != nil {
		return nil, err

	}

	err = json.Unmarshal(cfgMap["port"], &port)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfgMap["hipchat"], &hcClient)
	if err != nil {
		return nil, err
	}

	return cfgMap, nil
}
