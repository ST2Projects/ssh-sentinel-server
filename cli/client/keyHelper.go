package clientcli

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

func IsCertValid(certPath string) (bool, string) {

	certValid := PathExists(certPath)

	if certValid {
		certBytes, err := os.ReadFile(certPath)
		if err != nil {
			log.Errorf("%s - cert does not exist or cannot be read", certPath)
		}

		pub, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)
		if err != nil {
			log.Errorf("Error when parsing cert: %s", err.Error())
		}

		cert, ok := pub.(*ssh.Certificate)

		if !ok {
			log.Errorf("Failed to cast to cert")
		}

		now := time.Now().UTC()
		validBefore := time.Unix(int64(cert.ValidBefore), 0).UTC()
		validAfter := time.Unix(int64(cert.ValidAfter), 0).UTC()

		validBeforeString := validBefore.Format("2006-01-02 15:04:05.5 UTC")

		return now.After(validAfter) && now.Before(validBefore), validBeforeString
	}

	return false, "Cert not found"
}
