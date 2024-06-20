package server

import (
	"crypto/rand"
	"encoding/json"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/helper"
	"github.com/st2projects/ssh-sentinel-server/model/api"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

func KeySignHandler(context echo.Context) error {

	signRequest, err := marshallSigningRequest(context.Request().Body)

	if err != nil {
		return helper.NewError("Failed to marshall signing request: [%s]", err)
	}

	// ssh.ParseAuthorizedKey expects the key to be in the "disk" format
	pubKeyDisk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signRequest.Key))
	if err != nil {
		return helper.NewError("Failed to parse key: [%s]", err)
	}

	authRecord, _ := context.Get(apis.ContextAuthRecordKey).(*models.Record)

	cert, err := makeSSHCertificate(pubKeyDisk, authRecord.Username(), signRequest.Principals, signRequest.Extensions)
	if err != nil {
		return helper.NewError("Failed to sign cert: [%s]", err)
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	var response = api.NewKeySignResponse(true, "")
	response.SignedKey = string(signedCert)

	return context.JSON(200, response)
}

func marshallSigningRequest(requestReader io.Reader) (api.KeySignRequest, error) {

	body, err := io.ReadAll(requestReader)
	signRequest := api.KeySignRequest{}

	if err != nil {
		return api.KeySignRequest{}, err
	}

	json.Unmarshal(body, &signRequest)
	return signRequest, err
}

func makeSSHCertificate(pubKey ssh.PublicKey, username string, principals []string, extensions []api.Extension) (*ssh.Certificate, error) {
	caPriv := getCAKey()

	validBefore, validAfter := computeValidity()

	cert := &ssh.Certificate{
		Key:             pubKey,
		Serial:          0,
		CertType:        ssh.UserCert,
		ValidPrincipals: principals,
		ValidAfter:      validAfter,
		ValidBefore:     validBefore,
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{
				"source-address": "0.0.0.0/0",
			},
		},
	}

	cert.Extensions = getExtensionsAsMap(extensions, username)

	err := cert.SignCert(rand.Reader, caPriv)

	return cert, err
}

func getExtensionsAsMap(extensions []api.Extension, username string) map[string]string {
	if extensions == nil || len(extensions) == 0 {
		log.Warnf("No extensions found in request for user [%s]. Using default extensions %s", username, config.Config.DefaultExtensions)
		return mapExtensions(config.Config.DefaultExtensions)
	} else {
		log.Infof("User [%s] is adding %s extensions", username, extensions)
		return mapExtensions(extensions)
	}
}

func mapExtensions(extensions []api.Extension) map[string]string {
	mappedExtensions := map[string]string{}

	for _, extension := range extensions {
		mappedExtensions[string(extension)] = ""
	}

	return mappedExtensions
}

func computeValidity() (uint64, uint64) {
	now := time.Now()
	validBefore := uint64(now.Unix())
	maxDuration, _ := time.ParseDuration(config.Config.MaxValidTime)
	validAfter := uint64(now.Add(maxDuration).Unix())

	return validAfter, validBefore
}

func getCAKey() (caPriv ssh.Signer) {

	_, privKey, err := GetCAKeyPair()
	if err != nil {
		panic(err)
	}

	return privKey
}
