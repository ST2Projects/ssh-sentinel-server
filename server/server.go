package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-core/model"
	"github.com/st2projects/ssh-sentinel-server/config"
	"github.com/st2projects/ssh-sentinel-server/helper"
	cmdModel "github.com/st2projects/ssh-sentinel-server/model"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

	cert, err := MakeSSHCertificate(pubKeyDisk, signRequest.Username, signRequest.Principals, signRequest.Extensions)
	if err != nil {
		panic(helper.NewError("Failed to sign cert: [%s]", err))
	}

	signedCert := ssh.MarshalAuthorizedKey(cert)

	var response = model.NewKeySignResponse(true, "")
	response.SignedKey = string(signedCert)

	writer.WriteHeader(http.StatusOK)
	err = responseEncoder.Encode(response)

	if err != nil {
		panic(helper.NewError("Failed to encode response: [%s]", err))
	}
}

func MarshallSigningRequest(requestReader io.Reader) (model.KeySignRequest, error) {

	body, err := io.ReadAll(requestReader)
	signRequest := model.KeySignRequest{}

	if err == nil {
		json.Unmarshal(body, &signRequest)
	}

	return signRequest, err
}

func MakeSSHCertificate(pubKey ssh.PublicKey, username string, principals []string, extensions []model.Extension) (*ssh.Certificate, error) {
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

func getExtensionsAsMap(extensions []model.Extension, username string) map[string]string {
	if extensions == nil || len(extensions) == 0 {
		log.Warnf("No extensions found in request for user [%s]. Using default extensions %s", username, config.Config.DefaultExtensions)
		return mapExtensions(config.Config.DefaultExtensions)
	} else {
		log.Infof("User [%s] is adding %s extensions", username, extensions)
		return mapExtensions(extensions)
	}
}

func mapExtensions(extensions []model.Extension) map[string]string {
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

	keyFile := config.Config.CAPrivateKey
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

func Serve(httpConfig *cmdModel.HTTPConfig) {

	var (
		certReloader *simplecert.CertReloader
		err          error
		numRenews    int
		ctx, cancel  = context.WithCancel(context.Background())

		// init strict tlsConfig (this will enforce the use of modern TLS configurations)
		// you could use a less strict configuration if you have a customer facing web application that has visitors with old browsers
		tlsConf = tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)

		// a simple constructor for a http.Server with our Handler
		makeServer = func() *http.Server {
			return &http.Server{
				Addr:      fmt.Sprintf("0.0.0.0:%d", httpConfig.HttpsPort),
				Handler:   makeRouter(),
				TLSConfig: tlsConf,
			}
		}

		// init server
		srv = makeServer()

		// init simplecert configuration
		cfg = simplecert.Default
	)

	configuredTls := config.GetTLSConfig()
	cfg.Local = configuredTls.Local
	cfg.CacheDir = "./resources"
	cfg.Domains = configuredTls.CertDomains
	cfg.SSLEmail = configuredTls.CertEmail

	if configuredTls.DNSProvider != "" {
		cfg.DNSProvider = configuredTls.DNSProvider
	}

	if configuredTls.Local {
		cfg.HTTPAddress = "localhost"
		cfg.TLSAddress = "localhost"
	}

	cfg.WillRenewCertificate = func() {
		cancel()
	}

	cfg.DidRenewCertificate = func() {
		numRenews++
		// Restart the server
		ctx, cancel = context.WithCancel(context.Background())
		srv = makeServer()

		// Force reload the cert
		certReloader.ReloadNow()

		go serve(ctx, srv)
	}

	certReloader, err = simplecert.Init(cfg, func() {
		os.Exit(0)
	})

	if err != nil {
		log.Fatalf("Simple cert init failed: %s\n", err)
	}

	// Redirect 80 -> 443
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", httpConfig.HttpPort), http.HandlerFunc(simplecert.Redirect))

	tlsConf.GetCertificate = certReloader.GetCertificateFunc()
	log.Infof("Serving at https://%s:%d", configuredTls.CertDomains[0], httpConfig.HttpsPort)
	serve(ctx, srv)
	<-make(chan bool)
}

func makeRouter() *mux.Router {
	commonHandlers := alice.New(LoggingHandler, ErrorHandler)
	authHandlers := commonHandlers.Append(AuthenticationHandler)

	router := mux.NewRouter()

	router.Handle("/", commonHandlers.ThenFunc(PingHandler))
	router.Handle("/ping", commonHandlers.ThenFunc(PingHandler))
	router.Handle("/capubkey", commonHandlers.ThenFunc(CAPubKeyHandler))
	router.Handle("/ssh", authHandlers.ThenFunc(KeySignHandler))
	return router
}

func serve(ctx context.Context, srv *http.Server) {
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen %s\n", err)
		}
	}()

	log.Info("Server started")
	<-ctx.Done()
	log.Info("Server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5+time.Second)
	defer func() {
		cancel()
	}()

	err := srv.Shutdown(ctxShutdown)
	if err == http.ErrServerClosed {
		log.Info("Server stopped correctly")
	} else {
		log.Errorf("Error when stopping server %s\n", err)
	}
}
