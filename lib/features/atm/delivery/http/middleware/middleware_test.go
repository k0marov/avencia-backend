package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/middleware"
)

func TestATMAuthMiddleware(t *testing.T) {
	authHeader := RandomString()

	t.Run("error case - authorization header is not right", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("", "/", nil)
		request.Header.Set("Authorization", authHeader)

		nextCalled := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
		})

		validator := func([]byte) error {
			return client_errors.InvalidATMSecret
		}

		middleware.NewATMAuthMiddleware(validator)(nextHandler).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.InvalidATMSecret)
		Assert(t, nextCalled, false, "next handler was called")
	})

	t.Run("happy case", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("", "/", nil)
		request.Header.Set("Authorization", authHeader)

		nextCalled := false 
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
		})

		validator := func(gotAtmSecret []byte) error {
			if string(gotAtmSecret) == authHeader {
				return nil
			}
			panic("unexpected")
		}

		middleware.NewATMAuthMiddleware(validator)(nextHandler).ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
		Assert(t, nextCalled, true, "next handler was called")
	})
}
