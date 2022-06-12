package src

import (
	"crypto/rand"
	"encoding/json"
	"github.com/justinas/alice"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"ssh-sentinel-server/src/model"
	"time"
)

type KeySignRequest struct {
	APIKey     string   `json:"api_key"`
	Principals []string `json:"principals"`
	Key        string   `json:"key"`
}

const contentTypeKey = "Content-Type"
const jsonContentType = "application/json"

func KeySignHandler(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set(contentTypeKey, jsonContentType)
	responseEncoder := json.NewEncoder(writer)

	signRequest, err := MarshallSigningRequest(request)

	if err != nil {
		HandleError(err, "Failed to marshall signing request", responseEncoder)
	}

	// ssh.ParseAuthorizedKey expects the key to be in the "disk" format
	pubKeyDisk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signRequest.Key))

	if err != nil {
		HandleError(err, "Failed to parse key from request", responseEncoder)
	}

	cert, err := MakeSSHCertificate(pubKeyDisk, nil, 0, 0)

	if err != nil {
		HandleError(err, "Failed to sign key", responseEncoder)
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	var response = model.NewKeySignResponse(true, "")
	response.SignedKey = string(signedCert)

	writer.WriteHeader(http.StatusOK)
	responseEncoder.Encode(response)

}

func HandleError(err error, msg string, responseEncoder *json.Encoder) {

	response := model.NewKeySignResponse(false, msg+" "+err.Error())

	responseEncoder.Encode(response)
}

func MarshallSigningRequest(request *http.Request) (KeySignRequest, error) {

	body, err := ioutil.ReadAll(request.Body)
	signRequest := KeySignRequest{}

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

	privKeyFile, err := ioutil.ReadFile("resources/CA")

	if err != nil {
		log.Fatal("Cannot load priv key", err)
	}

	privKey, err := ssh.ParsePrivateKey(privKeyFile)

	if err != nil {
		log.Fatal("Cannot parse privKey", err)
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
				response := model.NewKeySignResponse(false, err)
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

func Serve() {

	commonHandlers := alice.New(LoggingHandler, ErrorHandler)
	mux := http.NewServeMux()

	mux.HandleFunc("/", Version)
	mux.HandleFunc("/version", Version)
	mux.Handle("/ssh", commonHandlers.ThenFunc(KeySignHandler))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
