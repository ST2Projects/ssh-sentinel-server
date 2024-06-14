package server

import (
	"github.com/st2projects/ssh-sentinel-server/model/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPrincipalsReturnsTrueWhenRequestedAndAllowed(t *testing.T) {

	allowed := []db.Principal{db.AsPrincipal("test")}

	result := CheckPrincipals(allowed, []string{"test"})

	assert.True(t, result)
}

func TestCheckPrincipalsReturnsFalseWhenNotAllowed(t *testing.T) {

	result := CheckPrincipals([]db.Principal{}, []string{"test"})

	assert.False(t, result)
}
