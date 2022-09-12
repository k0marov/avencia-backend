package http_helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	authService "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("contentType", "application/json")
}

func WriteJson(w http.ResponseWriter, obj any) {
	setJsonHeader(w)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		ThrowHTTPError(w, err)
		return
	}
}

func GetUserOrAddUnauthorized(w http.ResponseWriter, r *http.Request) (authEntities.User, bool) {
	authUser, err := authService.UserFromCtx(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return authEntities.User{}, false
	}
	return authUser, true
}

func ThrowHTTPError(w http.ResponseWriter, err error) {
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
