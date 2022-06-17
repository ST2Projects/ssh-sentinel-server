package sql

import (
	"github.com/google/uuid"
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

	dbConnection.AutoMigrate(&db.User{}, &db.Principal{})
	dbConnection.Model(&db.User{}).Association("user_id")

	p := db.AsPrincipal("abc")
	var all []db.Principal
	all = append(all, p)
	all = append(all, db.AsPrincipal("def"))

	dbConnection.Create(&db.User{

		ID:         0,
		Name:       "Ben",
		Expired:    false,
		Principals: all,
	})

	user := GetUserByAPIKey(uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	log.Info(user)
	log.Infof("Name: [%s], ID [%d], Expired [%t], Principals [%v]", user.Name, user.ID, user.Expired, GetPrincipalsByUser(user))
}

func GetUserByAPIKey(apiKey uuid.UUID) db.User {
	var user db.User
	dbConnection.First(&user, "api_key = ?", apiKey)
	return user
}

func GetPrincipalsByUser(user db.User) []db.Principal {
	var principals []db.Principal
	dbConnection.Find(&principals, "user_id = ?", user.ID)
	return principals
}
