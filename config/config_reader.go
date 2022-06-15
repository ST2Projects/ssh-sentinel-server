package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	CAPrivateKey string `json:"CAPrivateKey"`
	CAPublicKey  string `json:"CAPublicKey"`
	MaxValidTime string `json:"MaxValidTime"`
}

func NewConfig(configFile string) Config {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		panic("config file " + configFile + " does not exist")
	}
	configString, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(err)
	}

	appConfig := Config{}

	json.Unmarshal(configString, &appConfig)

	return appConfig
}
