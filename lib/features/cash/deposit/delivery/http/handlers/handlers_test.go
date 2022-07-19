package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	handlers2 "github.com/k0marov/avencia-backend/lib/features/cash/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/responses"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGenerateCodeHandler(t *testing.T) {
	http_test_helpers.BaseTest401(t, handlers2.NewGenerateCodeHandler(nil))

	user := RandomUser()
	requestWithUser := http_test_helpers.AddAuthDataToRequest(http_test_helpers.CreateRequest(nil), user)

	t.Run("happy case", func(t *testing.T) {
		response := httptest.NewRecorder()
		code := RandomString()
		expAt := time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC)
		generate := func(gotUser auth.User) (string, time.Time, error) {
			if gotUser == user {
				return code, expAt, nil
			}
			panic("unexpected")
		}
		handlers2.NewGenerateCodeHandler(generate)(response, requestWithUser)

		AssertJSONData(t, response, responses.CodeResponse{TransactionCode: code, ExpiresAt: expAt.Unix()})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, w *httptest.ResponseRecorder) {
		generate := func(auth.User) (string, time.Time, error) {
			return "", time.Time{}, err
		}
		handlers2.NewGenerateCodeHandler(generate)(w, requestWithUser)
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
		handlers2.NewVerifyCodeHandler(verify)(response, request)

		AssertJSONData(t, response, responses.VerifiedCodeResponse{UserInfo: responses.UserInfoResponse{Id: userInfo.Id}})
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		verify := func(string) (entities.UserInfo, error) {
			return entities.UserInfo{}, err
		}
		handlers2.NewVerifyCodeHandler(verify)(response, request)
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

		handlers2.NewCheckBanknoteHandler(checker)(response, request)

		AssertJSONData(t, response, responses.AcceptionResponse{Accept: accept})
	})
}

func TestFinalizeTransactionHandler(t *testing.T) {
	t.Run("should call service and return the result in the 'accept' field", func(t *testing.T) {
		transaction := RandomTransactionData()
		transactionJson, _ := json.Marshal(handlers.FinalizeTransactionRequest{
			UserId:    transaction.UserId,
			ATMSecret: string(transaction.ATMSecret),
			Currency:  transaction.Currency,
			Amount:    transaction.Amount,
		})
		accept := RandomBool()

		request := http_test_helpers.CreateRequest(bytes.NewReader(transactionJson))
		response := httptest.NewRecorder()

		finalizer := func(trans values.TransactionData) bool {
			if reflect.DeepEqual(trans, transaction) {
				return accept
			}
			panic("unexpected")
		}

		handlers2.NewFinalizeTransactionHandler(finalizer)(response, request)

		AssertJSONData(t, response, responses.AcceptionResponse{Accept: accept})
	})
}
