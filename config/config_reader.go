package config

import (
	"encoding/json"
	"errors"
	"github.com/st2projects/ssh-sentinel-server/model/api"
	"os"
)

type ConfigType struct {
	MaxValidTime      string          `json:"maxValidTime"`
	DefaultExtensions []api.Extension `json:"defaultExtensions"`
}

var Config *ConfigType

func MakeConfig(configPath string) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic("config file " + configPath + " does not exist")
	}

	configString, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	appConfig := ConfigType{}

	json.Unmarshal(configString, &appConfig)
	Config = &appConfig
}
