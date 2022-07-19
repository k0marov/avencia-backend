package lib

import (
	"context"
	"github.com/k0marov/avencia-backend/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit"
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
	authMiddleware := auth.NewAuthMiddleware(fbApp)

	cashDepositHandlers := deposit.NewCashDepositHandlers(conf)

	return api.NewAPIRouter(cashDepositHandlers, authMiddleware)
}
