package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	_ "github.com/st2projects/ssh-sentinel-server/migrations"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

const (
	DataDir   = "./pb_data"
	PublicDir = "./pb_public"
)

var AppContext *pocketbase.PocketBase

func Start() {
	AppContext = pocketbase.New()

	onAfterBootstrap(AppContext)
	onBeforeServe(AppContext)

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	log.Infof("AutoMigration [%v] arg0 [%s] tmpDir [%s]", isGoRun, os.Args[0], os.TempDir())

	migratecmd.MustRegister(AppContext, AppContext.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	if err := AppContext.Start(); err != nil {
		log.Fatal(err)
	}
}

func onAfterBootstrap(app *pocketbase.PocketBase) {
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		config.MakeConfig(DataDir + "/config.json")
		log.Infof("Read app config %+v", config.Config)

		_, err := app.Dao().FindFirstRecordByData("caKeys", "default", true)

		if err != nil {
			log.Warnf("No default key found, creating one")
			err = createNewDefaultKey(app)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func onBeforeServe(app *pocketbase.PocketBase) {

	commonHandlers := []echo.MiddlewareFunc{apis.ActivityLogger(app)}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/ui", apis.StaticDirectoryHandler(os.DirFS(PublicDir), true))
		e.Router.GET("/ping", PingHandler, append(commonHandlers)...) // TODO remove auth
		e.Router.GET("/capubkey", CAPubKeyHandler, commonHandlers...)
		e.Router.POST("/sign", KeySignHandler, append(commonHandlers, apis.RequireRecordAuth("users"))...)

		return nil
	})
}

func createNewDefaultKey(app *pocketbase.PocketBase) error {
	collection, err := app.Dao().FindCollectionByNameOrId("caKeys")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)

	privKey, pubKey := makeKeyPair()

	err = form.LoadData(map[string]any{
		"pubKey":  pubKey,
		"privKey": privKey,
		"default": true,
	})
	if err != nil {
		return err
	}

	if err = form.Submit(); err != nil {
		return err
	}

	log.Infof("Created new default key pubkey: %s", pubKey)

	return nil
}

func makeKeyPair() (string, string) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	//
	//signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		panic(err) // This is throwing
	}

	marshalledPrivateKey, err := ssh.MarshalPrivateKey(privateKey, "SSH CA")
	if err != nil {
		panic(err)
	}

	writer := new(strings.Builder)
	err = pem.Encode(writer, marshalledPrivateKey)
	if err != nil {
		panic(err)
	}

	pemPrivateKeyString := writer.String()
	writer.Reset()

	signer, err := ssh.ParsePrivateKey([]byte(pemPrivateKeyString))
	if err != nil {
		panic(err)
	}

	publicKeyString := string(ssh.MarshalAuthorizedKey(signer.PublicKey()))

	return pemPrivateKeyString, publicKeyString
}
