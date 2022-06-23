package sql

import (
	log "github.com/sirupsen/logrus"
	_ "gorm.io/driver/sqlite" // Import sqlite3 driver
	"gorm.io/gorm"
	"ssh-sentinel-server/config"
	"ssh-sentinel-server/model/db"
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

	err := dbConnection.AutoMigrate(&db.User{}, &db.Principal{}, &db.APIKeyType{})

	if err != nil {
		log.Fatalf("Failed to perform migration: [%s]", err.Error())
	}

	dbConnection.Model(&db.User{}).Association("user_id")
	dbConnection.Model(&db.APIKeyType{}).Association("apiKey")
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

	var apiKey db.APIKeyType
	dbConnection.Find(&apiKey, "user_id = ?", user.ID)
	user.APIKey = apiKey

	return user
}

func GetUserByID(id uint) db.User {
	var user = db.User{}
	dbConnection.First(&user, "id = ?", id)
	return user
}

func GetAPIKey(apiKey db.APIKeyType) db.APIKeyType {
	var key db.APIKeyType
	dbConnection.First(&key, "key = ?", apiKey.Key)
	return key
}

func GetApiKeyByID(id uint) db.APIKeyType {
	var key db.APIKeyType
	dbConnection.First(&key, "user_id = ?", id)
	return key
}

func GetPrincipalsByID(id uint) []db.Principal {
	var principals []db.Principal
	dbConnection.Find(&principals, "user_id = ?", id)
	return principals
}
