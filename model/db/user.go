package db

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/crypto"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         int         `db:"id" json:"id"`
	UserName   string      `db:"username" json:"user_name" gorm:"unique"`
	Name       string      `db:"name" json:"name"`
	Expired    bool        `db:"expired" json:"expired"`
	APIKey     APIKey      `db:"keyId" json:"keyId"`
	Principals []Principal `db:"principals" json:"principals"`
}

type Principal struct {
	gorm.Model
	UserID    uint   `db:"id" json:"id"`
	Principal string `db:"principal" json:"principal"`
}

type APIKey struct {
	gorm.Model
	UserId uint   `db:"id" json:"id"`
	Key    string `json:"api_key" db:"apiKey"`
}

func NewAPIKey() (APIKey, string) {
	id := uuid.New().String()

	apiKey := APIKey{}

	k, err := crypto.GenerateHash(crypto.PasswordConfig{}.DefaultConfig(), id)

	if err != nil {
		log.Fatal("Cannot create key", err)
	}

	apiKey.Key = k

	return apiKey, id
}

func (p *Principal) String() string {
	return p.Principal
}

func AsPrincipal(name string) Principal {
	return Principal{Principal: name}
}

func NewUser(name string, username string, principals []string) (User, string) {

	var principalTypes []Principal

	for _, p := range principals {
		principalTypes = append(principalTypes, AsPrincipal(p))
	}

	apiKey, rawKey := NewAPIKey()

	return User{
		Name:       name,
		UserName:   username,
		Expired:    false,
		Principals: principalTypes,
		APIKey:     apiKey,
	}, rawKey
}

func (u *User) Table() string {
	return "user"
}

func (u *User) IDColumn() string {
	return "id"
}
