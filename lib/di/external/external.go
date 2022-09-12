package external

import (
	"context"
	"io/ioutil"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/AvenciaLab/avencia-backend/lib/config"
	"github.com/AvenciaLab/avencia-backend/lib/core/db/foundationdb"
	"github.com/AvenciaLab/avencia-backend/lib/di"
	authStoreImpl "github.com/AvenciaLab/avencia-backend/lib/features/auth/store"
	"google.golang.org/api/option"
)

func initFirebase(config config.Config) *firebase.App {
	opt := option.WithCredentialsFile(config.FirebaseSecretPath)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	return fbApp
}
func InitializeExternal() di.ExternalDeps {
	conf := config.LoadConfig()
	// ===== CONFIG =====
	atmSecret, err := ioutil.ReadFile(conf.ATMSecretPath)
	if err != nil {
		log.Fatalf("error while reading atm secret: %v", err)
	}
	jwtSecret, err := ioutil.ReadFile(conf.JWTSecretPath)
	if err != nil {
		log.Fatalf("error while reading jwt secret: %v", err)
	}

	// ===== FIREBASE =====
	fbApp := initFirebase(conf)
	fbAuth, err := fbApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error while initializing firebase auth: %v", err)
	}

	// ===== DB =====
	fdb.MustAPIVersion(710)
	foundationDB := fdb.MustOpenDefault()

	runTrans := foundationdb.NewTransactionRunner(foundationDB)

	authFacade := authStoreImpl.NewFBAuthFacade(fbAuth)

	return di.ExternalDeps{
		AtmSecret: atmSecret,
		JwtSecret: jwtSecret,
		Auth:      authFacade,
		TRunner:   runTrans,
	}

}
