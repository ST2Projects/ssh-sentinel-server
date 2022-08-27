package db

import "testing"

func TestUser_Table(t *testing.T) {
	u := User{}

	tableName := u.Table()

	if !(tableName == "user") {
		t.Errorf("Incorrect table name %s . Expected 'user'", tableName)
	}
}
