package clientcli

import (
	"encoding/json"
	"github.com/st2projects/ssh-sentinel-server/model/api"
	"os"
)

type ClientConfigType struct {
	EndPoint   string          `json:"endPoint"`
	APIKey     string          `json:"apiKey"`
	Username   string          `json:"username"`
	Principals []string        `json:"principals"`
	Extensions []api.Extension `json:"extensions"`
	PublicKey  string          `json:"publicKey"`
	CertFile   string          `json:"certFile"`
}

var Config *ClientConfigType

func MakeConfig(configFile string) {
	if !PathExists(configFile) {
		panic("config file " + configFile + " does not exits")
	}

	configString, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configString, &Config)
	if err != nil {
		panic(err)
	}

}

func (c *ClientConfigType) GetPublicKey() string {
	return ExpandPath(c.PublicKey)
}

func (c *ClientConfigType) GetCertFile() string {
	return ExpandPath(c.CertFile)
}
