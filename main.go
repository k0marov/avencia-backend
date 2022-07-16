package main

import (
	"context"
	"fmt"
	"github.com/k0marov/avencia-backend/auth"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("firebase_secret.json")
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	fbAuth, err := fbApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firebase Auth: %v", err)
	}

	authMiddleware := auth.FirebaseHttpMiddleware{AuthClient: fbAuth}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.UserFromCtx(r.Context())
		if err != nil {
			w.Write([]byte("You are not logged in."))
		} else {
			w.Write([]byte(fmt.Sprintf("You are logged in. The user id is: %s", user.Id)))
		}
	}

	http.ListenAndServe(":4242", authMiddleware.Middleware(http.HandlerFunc(testHandler)))
}
