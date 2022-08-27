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
	"io/ioutil"
	"net/http"
	"time"
)

func AuthenticationHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeKey, jsonContentType)

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(helper.NewError("Failed to marshall request %s", err))
		}

		signRequest, err := MarshallSigningRequest(bytes.NewReader(body))

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if err != nil {
			panic(helper.NewError("Failed to marshall request %s", err))
		}

		user := sql.GetUserByUsername(signRequest.Username)

		hasValidAPIKey, err := crypto.Validate(signRequest.APIKey, user.APIKey.Key)

		if !hasValidAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			panic(helper.NewError("Unauthorised key"))
		}

		hasValidPrincipals := CheckPrincipals(user.Principals, signRequest.Principals)

		if !hasValidPrincipals {
			panic(helper.NewError("One or more unauthorised principals requested %v", signRequest.Principals))
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
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
