package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-core/model"
	"github.com/st2projects/ssh-sentinel-server/crypto"
	"github.com/st2projects/ssh-sentinel-server/helper"
	"github.com/st2projects/ssh-sentinel-server/sql"
	"io"
	"net/http"
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
			authorisationFailed(w, "No such user %s", signRequest.Username)
		}

		hasValidAPIKey, err := crypto.Validate(signRequest.APIKey, user.APIKey.Key)

		if !hasValidAPIKey {
			authorisationFailed(w, "Invalid API key for user %s", signRequest.Username)
		}

		hasValidPrincipals := CheckPrincipals(user.Principals, signRequest.Principals)

		if !hasValidPrincipals {
			authorisationFailed(w, "One or more unauthorised principals requested %v", signRequest.Principals)
		}

		log.Infof("User %s is authenticated", signRequest.Username)

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
				response := model.NewKeySignResponse(false, errorMsg)
				json.NewEncoder(w).Encode(response)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
