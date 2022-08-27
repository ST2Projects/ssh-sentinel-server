package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/model"
	"github.com/st2projects/ssh-sentinel-server/server"
	"github.com/st2projects/ssh-sentinel-server/sql"
)

func InitialiseApp(configPath string, devMode bool, httpConfig *model.HTTPConfig) {

	customLogFormat := new(log.TextFormatter)
	customLogFormat.TimestampFormat = "2022-01-01 01:01:01.123"
	customLogFormat.FullTimestamp = true

	log.SetFormatter(customLogFormat)

	log.Info("Starting Sentinel service")
	config.MakeConfig(configPath, devMode)
	sql.Connect()

	server.Serve(httpConfig)
}
