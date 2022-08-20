package middleware

import (
	"net/http"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
)



func NewATMAuthMiddleware(validateSecret validators.ATMSecretValidator) core.Middleware {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := r.Header.Get("Authorization") 	
			err := validateSecret([]byte(secret)) 
			if err != nil {
				http_helpers.ThrowHTTPError(w, err)
				return 
			} else {
				next.ServeHTTP(w, r)
			}
    })
  }
}
