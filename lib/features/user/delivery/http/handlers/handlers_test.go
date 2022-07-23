package handlers_test

import (
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/user/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfoHandler(t *testing.T) {

	http_test_helpers.BaseTest401(t, handlers.NewGetUserInfoHandler(nil))

	user := RandomUser()

	requestWithUser := http_test_helpers.AddAuthDataToRequest(http_test_helpers.CreateRequest(nil), user)

	t.Run("happy case", func(t *testing.T) {
		userInfo := RandomUserInfo()
		generate := func(userId string) (entities.UserInfo, error) {
			if userId == user.Id {
				return userInfo, nil
			}
			panic("unexpected")
		}

		response := httptest.NewRecorder()
		handlers.NewGetUserInfoHandler(generate)(response, requestWithUser)

		AssertJSONData(t, response, userInfo.ToResponse())
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		generate := func(string) (entities.UserInfo, error) {
			return entities.UserInfo{}, err
		}
		handlers.NewGetUserInfoHandler(generate)(w, requestWithUser)
	})
}
