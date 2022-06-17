package config

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
)

type Configtype struct {
	CAPrivateKey string `json:"CAPrivateKey"`
	CAPublicKey  string `json:"CAPublicKey"`
	MaxValidTime string `json:"MaxValidTime"`
	Db           DbType `json:"db"`
}

type DialectType string

func (d *DbType) AsGormConnection() gorm.Dialector {
	var dialect gorm.Dialector

	switch d.Dialect {
	case "sqlite3":
		if d.DevMode {
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
	DevMode    bool        `json:"devMode"`
}

type DbDriver string

var Config *Configtype

func MakeConfig(configFile string) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		panic("config file " + configFile + " does not exist")
	}
	configString, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(err)
	}

	appConfig := Configtype{}

	json.Unmarshal(configString, &appConfig)

	Config = &appConfig
}

func GetDBConfig() DbType {
	return Config.Db
}
