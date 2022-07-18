package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/responses"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/values"
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

		AssertJSONData(t, response, responses.CodeResponse{TransactionCode: code})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		generate := func(auth.User) (string, error) {
			return "", err
		}
		handlers.NewGenerateCodeHandler(generate)(w, requestWithUser)
	})
}

func TestVerifyCodeHandler(t *testing.T) {
	codeReq := handlers.CodeRequest{TransactionCode: RandomString()}
	codeReqBody, _ := json.Marshal(codeReq)
	request := http_test_helpers.CreateRequest(bytes.NewReader(codeReqBody))
	userInfo := RandomUserInfo()

	t.Run("happy case", func(t *testing.T) {
		verify := func(code string) (entities.UserInfo, error) {
			if code == codeReq.TransactionCode {
				return userInfo, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		handlers.NewVerifyCodeHandler(verify)(response, request)

		AssertJSONData(t, response, responses.UserInfoResponse{Id: userInfo.Id})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		verify := func(string) (entities.UserInfo, error) {
			return entities.UserInfo{}, err
		}
		handlers.NewVerifyCodeHandler(verify)(response, request)
	})

}

func TestCheckBanknoteHandler(t *testing.T) {
	t.Run("should call service and return the result in the 'accept' field", func(t *testing.T) {
		transactionCode := RandomString()
		currency := RandomString()
		amount := RandomInt()
		banknoteJson, _ := json.Marshal(handlers.BanknoteCheckRequest{
			TransactionCode: transactionCode,
			Currency:        currency,
			Amount:          amount,
		})
		request := http_test_helpers.CreateRequest(bytes.NewReader(banknoteJson))
		response := httptest.NewRecorder()

		accept := RandomBool()
		checker := func(code string, banknote values.Banknote) bool {
			if code == transactionCode && banknote.Amount == amount && banknote.Currency == currency {
				return accept
			}
			panic("unexpected")
		}

		handlers.NewCheckBanknoteHandler(checker)(response, request)

		AssertJSONData(t, response, responses.BanknoteCheckResponse{Accept: accept})
	})
}
