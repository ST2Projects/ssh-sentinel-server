package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"ssh-sentinel-server/config"
	"ssh-sentinel-server/helper"
	http2 "ssh-sentinel-server/model/http"
	"time"
)

const contentTypeKey = "Content-Type"
const jsonContentType = "application/json"

func KeySignHandler(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set(contentTypeKey, jsonContentType)
	responseEncoder := json.NewEncoder(writer)

	signRequest, err := MarshallSigningRequest(request.Body)

	if err != nil {
		panic(helper.NewError("Failed to marshall signing request: [%s]", err))
	}

	// ssh.ParseAuthorizedKey expects the key to be in the "disk" format
	pubKeyDisk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signRequest.Key))

	if err != nil {
		panic(helper.NewError("Failed to parse key: [%s]", err))
	}

	cert, err := MakeSSHCertificate(pubKeyDisk, nil)

	if err != nil {
		panic(helper.NewError("Failed to sign cert: [%s]", err))
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	var response = http2.NewKeySignResponse(true, "")
	response.SignedKey = string(signedCert)

	writer.WriteHeader(http.StatusOK)
	responseEncoder.Encode(response)

}

func MarshallSigningRequest(requestReader io.Reader) (http2.KeySignRequest, error) {

	body, err := ioutil.ReadAll(requestReader)
	signRequest := http2.KeySignRequest{}

	if err == nil {
		json.Unmarshal(body, &signRequest)
	}

	return signRequest, err
}

func MakeSSHCertificate(pubKey ssh.PublicKey, principals []string) (*ssh.Certificate, error) {
	caPriv := GetCAKey()

	validBefore, validAfter := ComputeValidity()

	cert := &ssh.Certificate{
		Key:             pubKey,
		Serial:          0,
		CertType:        ssh.UserCert,
		ValidPrincipals: principals,
		ValidAfter:      validAfter,
		ValidBefore:     validBefore,
		Permissions:     ssh.Permissions{},
	}

	err := cert.SignCert(rand.Reader, caPriv)

	return cert, err
}

func ComputeValidity() (uint64, uint64) {
	now := time.Now()
	validBefore := uint64(now.Unix())
	maxDuration, _ := time.ParseDuration(config.Config.MaxValidTime)
	validAfter := uint64(now.Add(maxDuration).Unix())

	return validAfter, validBefore
}

func GetCAKey() (caPriv ssh.Signer) {

	work, _ := os.Getwd()
	log.Printf("Working dir [%s]", work)
	keyFile := "resources/CA"
	privKeyFile, err := ioutil.ReadFile(keyFile)

	if err != nil {
		panic(helper.NewError("Failed to read private key [%s] : [%s]", keyFile, err))
	}

	privKey, err := ssh.ParsePrivateKey(privKeyFile)

	if err != nil {
		panic(err)
	}

	return privKey
}

func Version(response http.ResponseWriter, r *http.Request) {
	io.WriteString(response, "Version 1")
}

func Serve(port int) {

	commonHandlers := alice.New(LoggingHandler, ErrorHandler, AuthenticationHandler)
	mux := http.NewServeMux()

	mux.HandleFunc("/", Version)
	mux.HandleFunc("/version", Version)
	mux.Handle("/ssh", commonHandlers.ThenFunc(KeySignHandler))

	bindAddr := fmt.Sprintf(":%d", port)

	err := http.ListenAndServe(bindAddr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
