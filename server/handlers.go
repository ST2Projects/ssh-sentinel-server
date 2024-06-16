package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/crypto"
	"github.com/st2projects/ssh-sentinel-server/helper"
	"github.com/st2projects/ssh-sentinel-server/model/api"
	"github.com/st2projects/ssh-sentinel-server/sql"
	"io"
	"net/http"
	"os"
	"time"
)

func AuthenticationHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeKey, jsonContentType)

		body, err := io.ReadAll(r.Body)

		if err != nil {
			panic(helper.NewError("Failed to marshall request %s", err))
		}

		signRequest, err := MarshallSigningRequest(bytes.NewReader(body))

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err != nil {
			panic(helper.NewError("Failed to marshall request %s", err))
		}

		user, err := sql.GetUserByUsername(signRequest.Username)

		if err != nil {
			authorisationFailed(w, "No such user %s", helper.Sanitize(signRequest.Username))
		}

		hasValidAPIKey, err := crypto.Validate(signRequest.APIKey, user.APIKey.Key)

		if !hasValidAPIKey {
			authorisationFailed(w, "Invalid API key for user %s", helper.Sanitize(signRequest.Username))
		}

		hasValidPrincipals := CheckPrincipals(user.Principals, signRequest.Principals)

		if !hasValidPrincipals {
			// Sanitize the principals for logging
			helper.SanitizeStringSlice(signRequest.Principals)
			authorisationFailed(w, "One or more unauthorised principals requested %v", signRequest.Principals)
		}

		log.Infof("User %s is authenticated", helper.Sanitize(signRequest.Username))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func authorisationFailed(w http.ResponseWriter, msg string, args ...any) {
	w.WriteHeader(http.StatusUnauthorized)

	log.Errorf(msg, args...)

	panic(helper.NewError("Authentication failed"))
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
				response := api.NewKeySignResponse(false, errorMsg)
				json.NewEncoder(w).Encode(response)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func PingHandler(writer http.ResponseWriter, _ *http.Request) {

	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(fmt.Sprintf("Pong\nTime now is %s", time.Now().Format("2006-01-02 15:04:05"))))
	if err != nil {
		log.Errorf("Failed to write ping response %s", err.Error())
	}
}

func CAPubKeyHandler(writer http.ResponseWriter, _ *http.Request) {
	conf := config.Config

	pubKeyPath := conf.CAPublicKey
	pubKey, err := os.ReadFile(pubKeyPath)

	if err != nil {
		log.Errorf("Failed to read pub key %s", err)
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(pubKey)

	if err != nil {
		log.Errorf("faile to write response: %s", err.Error())
	}
}
