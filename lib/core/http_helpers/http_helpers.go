package http_helpers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"log"
	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
}

func WriteJson(w http.ResponseWriter, obj any) {
	setJsonHeader(w)
	json.NewEncoder(w).Encode(obj)
}

func GetUserOrAddUnauthorized(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return auth.User{}, false
	}
	return authUser, true
}

func HandleServiceError(w http.ResponseWriter, err error) {
	clientError, isClientError := err.(client_errors.ClientError)
	if isClientError {
		ThrowClientError(w, clientError)
	} else {
		log.Printf("Error while serving request: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func ThrowClientError(w http.ResponseWriter, clientError client_errors.ClientError) {
	setJsonHeader(w)
	errorJson, _ := json.Marshal(clientError)
	http.Error(w, string(errorJson), clientError.HTTPCode)
}
