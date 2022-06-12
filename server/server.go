package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/justinas/alice"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"ssh-sentinel-server/helper"
	model2 "ssh-sentinel-server/model"
	"time"
)

const contentTypeKey = "Content-Type"
const jsonContentType = "application/json"

func KeySignHandler(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set(contentTypeKey, jsonContentType)
	responseEncoder := json.NewEncoder(writer)

	signRequest, err := MarshallSigningRequest(request)

	if err != nil {
		panic(helper.NewError("Failed to marshall signing request: [%s]", err))
	}

	// ssh.ParseAuthorizedKey expects the key to be in the "disk" format
	pubKeyDisk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signRequest.Key))

	if err != nil {
		panic(helper.NewError("Failed to parse key: [%s]", err))
	}

	cert, err := MakeSSHCertificate(pubKeyDisk, nil, 0, 0)

	if err != nil {
		panic(helper.NewError("Failed to sign cert: [%s]", err))
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	var response = model2.NewKeySignResponse(true, "")
	response.SignedKey = string(signedCert)

	writer.WriteHeader(http.StatusOK)
	responseEncoder.Encode(response)

}

func MarshallSigningRequest(request *http.Request) (model2.KeySignRequest, error) {

	body, err := ioutil.ReadAll(request.Body)
	signRequest := model2.KeySignRequest{}

	if err == nil {
		json.Unmarshal(body, &signRequest)
	}

	return signRequest, err
}

func MakeSSHCertificate(pubKey ssh.PublicKey, principals []string, validBefore uint64, validAfter uint64) (*ssh.Certificate, error) {
	caPriv := GetCAKey()

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

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func ErrorHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// There is surely a better way to do this
				errorMsg := fmt.Sprintf("%s", err)
				response := model2.NewKeySignResponse(false, errorMsg)
				json.NewEncoder(w).Encode(response)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Version(response http.ResponseWriter, r *http.Request) {
	io.WriteString(response, "Version 1")
}

func Serve(port int) {

	commonHandlers := alice.New(LoggingHandler, ErrorHandler)
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
