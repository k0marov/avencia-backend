package lib

import (
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core/constants"
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

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.UserFromCtx(r.Context())
		if err != nil {
			w.Write([]byte("You are not logged in."))
		} else {
			w.Write([]byte(fmt.Sprintf("You are logged in. The user id is: %s", user.Id)))
		}
	})
	return authMiddleware(testHandler)
}
