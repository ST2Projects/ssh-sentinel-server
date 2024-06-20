package server

import (
	"golang.org/x/crypto/ssh"
)

func GetCAKeyPair() (ssh.PublicKey, ssh.Signer, error) {
	keyRecord, err := AppContext.Dao().FindFirstRecordByData("caKeys", "default", true)

	if err != nil {
		return nil, nil, err
	}

	privateKey, err := ssh.ParsePrivateKey([]byte(keyRecord.GetString("privKey")))
	if err != nil {
		return nil, nil, err
	}

	return privateKey.PublicKey(), privateKey, err
}
