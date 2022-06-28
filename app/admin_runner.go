package app

import (
	log "github.com/sirupsen/logrus"
	"ssh-sentinel-server/config"
	"ssh-sentinel-server/model/db"
	"ssh-sentinel-server/sql"
)

func RunAdmin(configPath string, createUser bool, name string, username string, principals []string) {
	config.MakeConfig(configPath, false)

	user, apiKey := db.NewUser(name, username, principals)

	sql.Connect()
	sql.NewUser(&user)

	log.Infof("Created new user with APIKey [%s] . This API Key MUST be kept private", apiKey)
}
