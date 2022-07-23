package app

import (
	log "github.com/sirupsen/logrus"
	"ssh-sentinel-server/config"
	"ssh-sentinel-server/server"
	"ssh-sentinel-server/sql"
)

func InitialiseApp(configPath string, devMode bool) {

	customLogFormat := new(log.TextFormatter)
	customLogFormat.TimestampFormat = "2022-01-01 01:01:01.123"
	customLogFormat.FullTimestamp = true

	log.SetFormatter(customLogFormat)

	log.Info("Starting Sentinel service")
	config.MakeConfig(configPath, devMode)
	sql.Connect()

	server.Serve()
}
