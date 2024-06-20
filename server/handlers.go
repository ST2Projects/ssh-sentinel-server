package server

import (
	"fmt"
	"github.com/labstack/echo/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
)

func PingHandler(context echo.Context) error {

	return context.String(200, fmt.Sprintf("Pong\nTime now is %s", time.Now().Format("2006-01-02 15:04:05")))
}

func CAPubKeyHandler(context echo.Context) error {

	pubkey, _, err := GetCAKeyPair()
	if err != nil {
		log.Errorf("Failed to read pub key %s", err)
	}

	marshalledKey := string(ssh.MarshalAuthorizedKey(pubkey))

	return context.String(200, marshalledKey)
}
