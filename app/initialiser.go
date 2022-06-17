package app

import (
	log "github.com/sirupsen/logrus"
	"ssh-sentinel-server/config"
	"ssh-sentinel-server/server"
	"ssh-sentinel-server/sql"
)

func InitialiseApp(port int, configPath string) {

	log.Info("Starting Sentinel service")
	config.MakeConfig(configPath)
	sql.Connect()

	server.Serve(port)
	log.Info("Started")
}
