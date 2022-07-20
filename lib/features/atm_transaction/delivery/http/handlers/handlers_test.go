package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/k0marov/avencia-backend/api"
	"github.com/k0marov/avencia-backend/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
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
		generate := func(gotUser auth.User, transType service.TransactionType) (string, time.Time, error) {
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
		generate := func(auth.User, service.TransactionType) (string, time.Time, error) {
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
		verify := func(code string, transType service.TransactionType) (entities.UserInfo, error) {
			if code == codeReq.TransactionCode && transType == tType {
				return userInfo, nil
			}
			panic("unexpected args")
		}
		response := httptest.NewRecorder()
		handlers.NewVerifyCodeHandler(verify)(response, request)

		AssertJSONData(t, response, api.VerifiedCodeResponse{UserInfo: api.UserInfoResponse{Id: userInfo.Id}})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		verify := func(string, service.TransactionType) (entities.UserInfo, error) {
			return entities.UserInfo{}, err
		}
		handlers.NewVerifyCodeHandler(verify)(response, request)
	})

}

func TestCheckBanknoteHandler(t *testing.T) {
	transactionCode := RandomString()
	currency := RandomString()
	amount := RandomInt()
	banknoteRequest := api.BanknoteCheckRequest{
		TransactionCode: transactionCode,
		Currency:        currency,
		Amount:          amount,
	}
	banknoteJson, _ := json.Marshal(banknoteRequest)
	wantBanknoteValue := values.NewBanknote(banknoteRequest)
	request := http_test_helpers.CreateRequest(bytes.NewReader(banknoteJson))

	t.Run("should call service and return status code 200 if there is no error", func(t *testing.T) {
		response := httptest.NewRecorder()

		checker := func(code string, banknote values.Banknote) error {
			if code == transactionCode && banknote == wantBanknoteValue {
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
	transaction := RandomTransactionData()
	transactionJson, _ := json.Marshal(api.FinalizeTransactionRequest{
		UserId:    transaction.UserId,
		ATMSecret: string(transaction.ATMSecret),
		Currency:  transaction.Currency,
		Amount:    transaction.Amount,
	})
	request := http_test_helpers.CreateRequest(bytes.NewReader(transactionJson))

	t.Run("should call service and return status code 200 if there is no error", func(t *testing.T) {
		response := httptest.NewRecorder()

		finalizer := func(trans values.TransactionData) error {
			if reflect.DeepEqual(trans, transaction) {
				return nil
			}
			panic("unexpected")
		}

		handlers.NewFinalizeTransactionHandler(finalizer)(response, request)
		AssertStatusCode(t, response, http.StatusOK)
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		finalizer := func(values.TransactionData) error {
			return err
		}
		handlers.NewFinalizeTransactionHandler(finalizer)(response, request)
	})
}
