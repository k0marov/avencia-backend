package http_helpers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/core/client_errors"
	"net/http"
)

func ThrowClientError(w http.ResponseWriter, cErr client_errors.ClientError) {
	w.WriteHeader(cErr.HTTPCode)
	json.NewEncoder(w).Encode(cErr)
}
