package config

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/model/api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type ConfigType struct {
	DevMode           bool            `json:"devMode"`
	CAPrivateKey      string          `json:"CAPrivateKey"`
	CAPublicKey       string          `json:"CAPublicKey"`
	MaxValidTime      string          `json:"maxValidTime"`
	DefaultExtensions []api.Extension `json:"defaultExtensions"`
	Db                DbType          `json:"db"`
}

type DialectType string

func (d *DbType) AsGormConnection() gorm.Dialector {
	var dialect gorm.Dialector

	switch d.Dialect {
	case "sqlite3":
		if Config.DevMode {
			log.Warnf("Dev mode enabled. Deleting DB [%s]", d.Connection)
			err := os.Remove(d.Connection)
			if err != nil {
				log.Errorf("failed to delete DB %s", err.Error())
			}
		}
		dialect = sqlite.Open(d.Connection)
	default:
		log.Fatalf("Unkown dialect [%s]", d.Dialect)
	}

	return dialect
}

type DbType struct {
	Dialect    DialectType `json:"dialect"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	Connection string      `json:"connection"`
	DBName     string      `json:"dbName"`
}

type DbDriver string

var Config *ConfigType

func MakeConfig(configFile string, devMode bool) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		panic("config file " + configFile + " does not exist")
	}

	configString, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	appConfig := ConfigType{}

	json.Unmarshal(configString, &appConfig)
	appConfig.DevMode = devMode

	Config = &appConfig

}

func GetDBConfig() DbType {
	return Config.Db
}
