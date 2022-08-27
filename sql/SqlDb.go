package sql

import (
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/model/db"
	_ "gorm.io/driver/sqlite" // Import sqlite3 driver
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

func Connect() {
	dbConfig := config.GetDBConfig()

	connectionString := dbConfig.Connection

	dbConnection, _ = gorm.Open(dbConfig.AsGormConnection(), &gorm.Config{})
	log.Infof("Created connection to %s", connectionString)

	initTables()
}

func initTables() {

	err := dbConnection.AutoMigrate(&db.User{}, &db.Principal{}, &db.APIKey{})

	if err != nil {
		log.Fatalf("Failed to perform migration: [%s]", err.Error())
	}

	dbConnection.Model(&db.User{}).Association("user_id")
	dbConnection.Model(&db.APIKey{}).Association("apiKey")
}

func NewUser(user *db.User) {
	dbConnection.Create(user)
}

func GetUserByUsername(username string) db.User {

	var user = db.User{}
	dbConnection.First(&user, "user_name = ? ", username)

	var principals []db.Principal
	dbConnection.Find(&principals, "user_id = ?", user.ID)
	user.Principals = principals

	var apiKey db.APIKey
	dbConnection.Find(&apiKey, "user_id = ?", user.ID)
	user.APIKey = apiKey

	return user
}

func GetUserByID(id uint) db.User {
	var user = db.User{}
	dbConnection.First(&user, "id = ?", id)
	return user
}

func GetAPIKey(apiKey db.APIKey) db.APIKey {
	var key db.APIKey
	dbConnection.First(&key, "key = ?", apiKey.Key)
	return key
}

func GetApiKeyByID(id uint) db.APIKey {
	var key db.APIKey
	dbConnection.First(&key, "user_id = ?", id)
	return key
}

func GetPrincipalsByID(id uint) []db.Principal {
	var principals []db.Principal
	dbConnection.Find(&principals, "user_id = ?", id)
	return principals
}
