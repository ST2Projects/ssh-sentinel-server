package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         int         `db:"id" json:"id"`
	UserName   string      `db:"username" json:"user_name" gorm:"unique"`
	Name       string      `db:"name" json:"name"`
	Expired    bool        `db:"expired" json:"expired"`
	APIKey     APIKeyType  `db:"keyId" json:"keyId"`
	Principals []Principal `db:"principals" json:"principals"`
}

type Principal struct {
	gorm.Model
	UserID    uint   `db:"id" json:"id"`
	Principal string `db:"principal" json:"principal"`
}

type APIKeyType struct {
	gorm.Model
	UserId uint   `db:"id" json:"id"`
	Key    string `json:"api_key" db:"apiKey"`
}

func NewAPIKey() (APIKeyType, string) {
	id := uuid.New().String()

	sha := sha256.New()
	sha.Write([]byte(id))
	finalValue := sha.Sum(nil)

	apiKey := APIKeyType{}

	apiKey.Key = hex.EncodeToString(finalValue)

	return apiKey, id
}

func AsAPIKey(key uuid.UUID) APIKeyType {
	sha := sha256.New()
	sha.Write([]byte(key.String()))
	finalValue := sha.Sum(nil)

	apiKey := APIKeyType{}
	apiKey.Key = hex.EncodeToString(finalValue)
	return apiKey
}

func (k *APIKeyType) Validate(other string) bool {
	sha := sha256.New()
	sha.Write([]byte(other))
	otherSum := sha.Sum(nil)

	thisDecoded, err := hex.DecodeString(k.Key)

	if err != nil {
		log.Fatal("Failed to decode key", err)
	}

	log.Infof("This: [%s], Other: [%s] :: [%s] ", k.Key, hex.EncodeToString(otherSum), other)

	// ConstantTimeCompare returns 1 when equal...
	return subtle.ConstantTimeCompare(thisDecoded, otherSum) == 1
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
