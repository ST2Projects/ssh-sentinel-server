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
	"github.com/pocketbase/pocketbase/tools/list"
	log "github.com/sirupsen/logrus"
	"github.com/st2projects/ssh-sentinel-server/config"
	_ "github.com/st2projects/ssh-sentinel-server/migrations"
	"github.com/st2projects/ssh-sentinel-server/model/db"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"time"
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
	onTerminate(AppContext)

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	log.Infof("AutoMigration [%v] arg0 [%s] tmpDir [%s]", isGoRun, os.Args[0], os.TempDir())

	migratecmd.MustRegister(AppContext, AppContext.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	if err := AppContext.Start(); err != nil {
		log.Fatal(err)
	}
}

func onTerminate(app *pocketbase.PocketBase) {
	app.OnTerminate().Add(func(e *core.TerminateEvent) error {

		return nil
	})
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
	recordAuthHandlers := append(commonHandlers, requireRecordAuth(app, "users"))

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/ui", apis.StaticDirectoryHandler(os.DirFS(PublicDir), true))
		e.Router.GET("/ping", PingHandler, commonHandlers...)
		e.Router.GET("/capubkey", CAPubKeyHandler, append(recordAuthHandlers, eventFeedLogger(app, db.Fetch))...)
		e.Router.POST("/sign", KeySignHandler, append(recordAuthHandlers, eventFeedLogger(app, db.Sign))...)

		return nil
	})
}

func eventFeedLogger(app *pocketbase.PocketBase, event db.Event) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				return err
			}

			eventFeed, err := app.Dao().FindCollectionByNameOrId("eventFeed")
			if err != nil {
				return err
			}
			newEvent := models.NewRecord(eventFeed)
			newEventForm := forms.NewRecordUpsert(app, newEvent)

			userRecord, err := getAuthUserRecord(c)

			err = newEventForm.LoadData(map[string]any{
				"event":     event,
				"user":      userRecord.Id,
				"eventTime": time.Now().UTC(),
			})

			if err != nil {
				return err
			}

			if err = newEventForm.Submit(); err != nil {
				return nil
			}

			return nil
		}
	}
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

func requireRecordAuth(app *pocketbase.PocketBase, optCollectionNames ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userRecord, err := getAuthUserRecord(c)
			if err != nil {
				return err
			}

			// check record collection name
			if len(optCollectionNames) > 0 && !list.ExistInSlice(userRecord.Collection().Name, optCollectionNames) {
				return apis.NewForbiddenError("The authorized record model is not allowed to perform this action.", nil)
			}

			approved := userRecord.GetBool("adminApproved")
			verified := userRecord.Verified()

			if !approved || !verified {
				log.Infof("User [%s (%s) tried to access [%s] but verification status: verified[%v] approved[%v]",
					userRecord.Username(), userRecord.Email(), c.Request().RequestURI, verified, approved)
				return apis.NewForbiddenError("User is not verified or approved", nil)
			}

			eventFeedLogger(app, db.Login)

			return next(c)
		}
	}
}

func getAuthUserRecord(context echo.Context) (*models.Record, *apis.ApiError) {
	record, _ := context.Get(apis.ContextAuthRecordKey).(*models.Record)
	if record == nil {
		return nil, apis.NewUnauthorizedError("The request requires valid record authorization token to be set.", nil)
	}

	return record, nil
}
