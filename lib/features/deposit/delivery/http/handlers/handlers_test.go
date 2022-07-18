package handlers_test

import (
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/responses"
	"net/http/httptest"
	"testing"
)

func TestGenerateCodeHandler(t *testing.T) {
	http_test_helpers.BaseTest401(t, handlers.NewGenerateCodeHandler(nil))

	user := RandomUser()
	requestWithUser := http_test_helpers.AddAuthDataToRequest(http_test_helpers.CreateRequest(nil), user)

	t.Run("happy case", func(t *testing.T) {
		response := httptest.NewRecorder()
		code := RandomString()
		generate := func(gotUser auth.User) (string, error) {
			if gotUser == user {
				return code, nil
			}
			panic("unexpected")
		}
		handlers.NewGenerateCodeHandler(generate)(response, requestWithUser)

		AssertJSONData(t, response, responses.CodeResponse{Code: code})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		generate := func(auth.User) (string, error) {
			return "", err
		}
		handlers.NewGenerateCodeHandler(generate)(w, requestWithUser)
	})
}
