package src

import (
	"crypto/rand"
	"encoding/json"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type sshHandler struct {
}

type KeySignResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SignedKey string `json:"signedKey"`
	NotBefore uint64 `json:"notBefore"`
	NotAfter  uint64 `json:"notAfter"`
}

type KeySignRequest struct {
	APIKey     string   `json:"api_key"`
	Principals []string `json:"principals"`
	Key        string   `json:"key"`
}

const contentTypeKey = "Content-Type"
const jsonContentType = "application/json"

func (h *sshHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		response := &KeySignResponse{
			Success: false,
			Message: err.Error(),
		}

		json.NewEncoder(writer).Encode(response)
	}

	signRequest := KeySignRequest{}
	json.Unmarshal(body, &signRequest)

	pubKeyDisk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signRequest.Key))

	if err != nil {
		writer.WriteHeader(400)
		log.Fatal(err)
	}

	cert, err := MakeSSHCertificate(pubKeyDisk, nil, 0, 0)

	if err != nil {
		writer.WriteHeader(400)
		io.WriteString(writer, "Failed to sign key: "+err.Error())
		panic(err)
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	response := &KeySignResponse{
		Success:   true,
		Message:   "",
		SignedKey: string(signedCert),
		NotBefore: 0,
		NotAfter:  0,
	}

	writer.Header().Set(contentTypeKey, jsonContentType)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)

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

func Version(response http.ResponseWriter, r *http.Request) {
	io.WriteString(response, "Version 1")
}

func Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Version)
	mux.HandleFunc("/version", Version)
	mux.Handle("/ssh", &sshHandler{})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
