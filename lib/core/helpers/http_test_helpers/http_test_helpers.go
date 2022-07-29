package http_test_helpers

import (
	test_helpers2 "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func AddAuthDataToRequest(r *http.Request, user auth.User) *http.Request {
	ctx := auth.AddUserToCtx(user, r.Context())
	return r.WithContext(ctx)
}

// CreateRequest Since I am using go-chi router, handlers can be independent of urls and http methods
func CreateRequest(body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/handler-should-not-care", body)
}

func BaseTestServiceErrorHandling(t *testing.T, callErroringHandler func(error, *httptest.ResponseRecorder)) {
	t.Helper()
	t.Run("service throws a client error", func(t *testing.T) {
		clientError := test_helpers2.RandomClientError()
		response := httptest.NewRecorder()
		callErroringHandler(clientError, response)
		test_helpers2.AssertClientError(t, response, clientError)
	})
	t.Run("service throws an internal error", func(t *testing.T) {
		response := httptest.NewRecorder()
		callErroringHandler(test_helpers2.RandomError(), response)
		test_helpers2.AssertStatusCode(t, response, http.StatusInternalServerError)
	})
}

func BaseTest401(t *testing.T, handlerWithPanickingService http.Handler) {
	t.Helper()
	t.Run("should return 401 if authentication details are not provided via context (using auth middleware)", func(t *testing.T) {
		request := CreateRequest(nil)
		response := httptest.NewRecorder()
		handlerWithPanickingService.ServeHTTP(response, request)

		test_helpers2.AssertStatusCode(t, response, http.StatusUnauthorized)
	})
}
