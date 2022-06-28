package db

import "testing"

func TestAPIKeyType_Validate(t *testing.T) {
	myKey := APIKeyType{Key: "e072bebc2cf6191881f0a1af2af353e1ded499e77b9d05a0425a25c3fce90807"}

	validated := myKey.Validate("39d61458-7c4c-4e58-a79d-37f02a448ca9")

	if !validated {
		t.Error("Key did not validate")
	}
}

func TestUser_Table(t *testing.T) {
	u := User{}

	tableName := u.Table()

	if !(tableName == "user") {
		t.Errorf("Incorrect table name %s . Expected 'user'", tableName)
	}
}
