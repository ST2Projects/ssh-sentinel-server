package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/model/db"
	"github.com/st2projects/ssh-sentinel-server/sql"
)

func RunAdmin(configPath string, createUser bool, name string, username string, principals []string) {
	config.MakeConfig(configPath, false)

	user, apiKey := db.NewUser(name, username, principals)

	sql.Connect()
	sql.NewUser(&user)

	log.Infof("Created new user with APIKey [%s] . This API Key MUST be kept private", apiKey)
}
