package external

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/AvenciaLab/avencia-backend/lib/config"
	"github.com/AvenciaLab/avencia-backend/lib/core/db/foundationdb"
	"github.com/AvenciaLab/avencia-backend/lib/di"
	authStoreImpl "github.com/AvenciaLab/avencia-backend/lib/features/auth/store"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"google.golang.org/api/option"
)

func readSecret(filepath string) string {
	contents, err := os.ReadFile(filepath)
	if err != nil {
	  log.Fatalf("while reading the contents of the %v secret file: %w", filepath, err)	
	}
	return strings.TrimSpace(string(contents))
}

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
	atmSecret := readSecret(conf.ATMSecretPath)
	jwtSecret := readSecret(conf.JWTSecretPath)

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
