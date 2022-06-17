package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         int         `db:"id" json:"id"`
	Name       string      `db:"name" json:"name"`
	Expired    bool        `db:"expired" json:"expired"`
	APIKey     uuid.UUID   `db:"apiKey" json:"api_key"`
	Principals []Principal `db:"principals" json:"principals"`
}

type Principal struct {
	gorm.Model
	UserID    uint   `db:"id" json:"id"`
	Principal string `db:"principal" json:"principal"`
}

func (p *Principal) String() string {
	return p.Principal
}

func AsPrincipal(name string) Principal {
	return Principal{Principal: name}
}

func NewUser(name string, principals []Principal) User {
	return User{
		Name:       name,
		Expired:    false,
		Principals: principals,
	}
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) IDColumn() string {
	return "id"
}
