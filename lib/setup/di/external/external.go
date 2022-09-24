package external

import (
	"context"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/AvenciaLab/avencia-backend/lib/core/db/foundationdb"
	authStoreImpl "github.com/AvenciaLab/avencia-backend/lib/features/auth/store"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"google.golang.org/api/option"
)

func removeEOF(contents []byte) []byte {
	var eof = string([]byte{10})
	return []byte(strings.TrimSuffix(string(contents), eof))
}

func readSecret(filepath string) []byte {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("while reading the contents of the %v secret file: %v", filepath, err)
	}
	return removeEOF(contents)
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
	simpleDB := foundationdb.NewSimpleDB(runTrans)

	authFacade := authStoreImpl.NewFBAuthFacade(fbAuth)

	return di.ExternalDeps{
		AtmSecret: atmSecret,
		JwtSecret: jwtSecret,
		Auth:      authFacade,
		TRunner:   runTrans,
		SimpleDB:  simpleDB,
	}

}
