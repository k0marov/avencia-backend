package lib

import (
	"context"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction"
	"github.com/k0marov/avencia-backend/lib/features/wallet"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/store"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/k0marov/avencia-backend/lib/features/auth"
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

func Initialize() http.Handler {
	conf := config.LoadConfig()

	fbApp := initFirebase(conf)
	fsClient, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error while initializing firestore client: %v", err)
	}

	authMiddleware := auth.NewAuthMiddleware(fbApp)
	getWallet := wallet.NewWalletGetterImpl(fsClient)

	atmTransactionHandlers := atm_transaction.NewATMTransactionHandlers(conf, getWallet, walletStore.NewBalanceUpdater, fsClient)

	return api.NewAPIRouter(atmTransactionHandlers, authMiddleware)
}
