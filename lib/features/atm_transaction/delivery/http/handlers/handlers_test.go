package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGenerateCodeHandler(t *testing.T) {
	http_test_helpers.BaseTest401(t, handlers.NewGenerateCodeHandler(nil))

	user := RandomUser()

	t.Run("error case - transaction_type is not provided", func(t *testing.T) {
		request := http_test_helpers.AddAuthDataToRequest(http_test_helpers.CreateRequest(nil), user)
		response := httptest.NewRecorder()

		handlers.NewGenerateCodeHandler(nil)(response, request)
		AssertClientError(t, response, client_errors.TransactionTypeNotProvided)
	})

	tType := RandomTransactionType()
	endpoint := fmt.Sprintf("/asdf?transaction_type=%s", tType)
	requestWithUser := http_test_helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodGet, endpoint, nil), user)
	t.Run("happy case", func(t *testing.T) {
		code := RandomString()
		expAt := time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
		generate := func(gotUser auth.User, transType values.TransactionType) (string, time.Time, error) {
			if gotUser == user && transType == tType {
				return code, expAt, nil
			}
			panic("unexpected")
		}

		response := httptest.NewRecorder()
		handlers.NewGenerateCodeHandler(generate)(response, requestWithUser)

		AssertJSONData(t, response, api.CodeResponse{TransactionCode: code, ExpiresAt: expAt.Unix()})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		generate := func(auth.User, values.TransactionType) (string, time.Time, error) {
			return "", time.Time{}, err
		}
		handlers.NewGenerateCodeHandler(generate)(w, requestWithUser)
	})
}

func TestVerifyCodeHandler(t *testing.T) {

	t.Run("error case - transaction type is not provided", func(t *testing.T) {
		request := http_test_helpers.CreateRequest(nil)
		response := httptest.NewRecorder()

		handlers.NewVerifyCodeHandler(nil)(response, request)

		AssertClientError(t, response, client_errors.TransactionTypeNotProvided)
	})

	codeReq := api.CodeRequest{TransactionCode: RandomString()}
	codeReqBody, _ := json.Marshal(codeReq)
	tType := RandomTransactionType()
	endpoint := fmt.Sprintf("/any?transaction_type=%s", tType)
	request := httptest.NewRequest(http.MethodOptions, endpoint, bytes.NewReader(codeReqBody))

	t.Run("happy case", func(t *testing.T) {
		userInfo := RandomUserInfo()
		verify := func(code string, transType values.TransactionType) (entities.UserInfo, error) {
			if code == codeReq.TransactionCode && transType == tType {
				return userInfo, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		handlers.NewVerifyCodeHandler(verify)(response, request)

		AssertJSONData(t, response, api.VerifiedCodeResponse{UserInfo: userInfo.ToResponse()})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		verify := func(string, values.TransactionType) (entities.UserInfo, error) {
			return entities.UserInfo{}, err
		}
		handlers.NewVerifyCodeHandler(verify)(response, request)
	})

}

func TestCheckBanknoteHandler(t *testing.T) {
	req := RandomBanknoteCheckRequest()
	banknoteJson, _ := json.Marshal(req)
	wantBanknoteValue := values.NewBanknote(req)
	request := http_test_helpers.CreateRequest(bytes.NewReader(banknoteJson))

	t.Run("should call service and return status code 200 if there is no error", func(t *testing.T) {
		response := httptest.NewRecorder()

		checker := func(code string, banknote values.Banknote) error {
			if code == req.TransactionCode && reflect.DeepEqual(banknote, wantBanknoteValue) {
				return nil
			}
			panic("unexpected")
		}

		handlers.NewCheckBanknoteHandler(checker)(response, request)
		AssertStatusCode(t, response, http.StatusOK)
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		checker := func(string, values.Banknote) error {
			return err
		}
		handlers.NewCheckBanknoteHandler(checker)(w, request)
	})
}

func TestFinalizeTransactionHandler(t *testing.T) {
	req := RandomFinalizeTransationRequest()
	transactionJson, _ := json.Marshal(req)
	transaction := values.NewTransactionData(req)
	request := http_test_helpers.CreateRequest(bytes.NewReader(transactionJson))

	t.Run("should call service and return status code 200 if there is no error", func(t *testing.T) {
		response := httptest.NewRecorder()

		finalizer := func(secret []byte, trans values.Transaction) error {
			if string(secret) == req.ATMSecret && reflect.DeepEqual(trans, transaction) {
				return nil
			}
			panic("unexpected")
		}

		handlers.NewFinalizeTransactionHandler(finalizer)(response, request)
		AssertStatusCode(t, response, http.StatusOK)
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		finalizer := func([]byte, values.Transaction) error {
			return err
		}
		handlers.NewFinalizeTransactionHandler(finalizer)(response, request)
	})
}
