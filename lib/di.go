package lib

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core/constants"
	"github.com/k0marov/avencia-backend/lib/features/deposit"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"google.golang.org/api/option"
)

func initFirebase() *firebase.App {
	opt := option.WithCredentialsFile(constants.FirebaseSecretPath)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	return fbApp
}

func Initialize() http.Handler {
	fbApp := initFirebase()
	authMiddleware := auth.NewAuthMiddleware(fbApp)

	r := chi.NewRouter()

	r.Route("/deposit", deposit.NewDepositRouterImpl(authMiddleware))

	return r
}
