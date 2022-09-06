package integration_test

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/di"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

// TODO: maybe get rid of separation between atm auth and user auth and just use atm secret as a special auth token
// TODO: maybe stop testing the handlers layer in here, but just call services directly 

const (
	AtmAuthSecret = "atm_test"
	JwtSecret     = "jwt_test"
)

func TestIntegration(t *testing.T) {
	users := []MockUser{
		{
			Token: RandomString(),
			Id: "sam",
			Email: "sam@skomarov.com",
		},
		{
			RandomString(),
			"john",
			"test@example.com",
		},
		{
			RandomString(),
			"bill",
			"test2@example.com",
		},
	}
	extDeps, cancelDBTrans := prepareExternalDeps(t, users)
	defer cancelDBTrans()
	apiDeps := di.InitializeBusiness(extDeps)

	encodeRequest := func(req any) *http.Request {
		body, err := json.Marshal(req)
		AssertNoError(t, err)
		return httptest.NewRequest("", "/asdf", bytes.NewReader(body))
	}
	addAuth := func(req *http.Request, authHeader string) *http.Request {
		req.Header.Add("Authorication", authHeader)
		return req
	}

	assertBalance := func(user MockUser, balance core.Money) {
		request := addAuth(encodeRequest(nil), user.Token)
		response := httptest.NewRecorder()
		handler := apiDeps.AuthMW(apiDeps.Handlers.App.GetUserInfo)
		handler.ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK) 
		var userInfo api.UserInfoResponse
		err := json.Unmarshal(response.Body.Bytes(), &userInfo)
		AssertNoError(t, err)
		Assert(t, userInfo.Wallet[string(balance.Currency)], balance.Amount.Num(), "balance")
	}
	generateQRCode := func(user MockUser, transType values.TransactionType) api.GenTransCodeResponse {
		reqBody := api.GenTransCodeRequest{
			TransactionType: string(transType),
		}
		request := addAuth(encodeRequest(reqBody), user.Token)
		response := httptest.NewRecorder() 
		handler := apiDeps.AuthMW(apiDeps.Handlers.App.GenCode)
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
		var codeResp api.GenTransCodeResponse 
		err := json.Unmarshal(response.Body.Bytes(), &codeResp)
		AssertNoError(t, err)
		return codeResp
	}
	startTrans := func(user MockUser, tType values.TransactionType, qrText string) (tId string) {
		reqBody := api.OnTransactionCreateRequest{
			TransactionReference: "asdf",
			TerminalId:           "1234",
			TerminalSid:          "4321",
			Type:                 string(tType),
			QRCodeText:           qrText,
		}
		response := httptest.NewRecorder()
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret)
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.OnCreate)
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
		var jsonResp api.OnTransactionCreateResponse 
		err := json.Unmarshal(response.Body.Bytes(), &jsonResp) 
		AssertNoError(t, err)
		return jsonResp.TransactionId
	}
	insertBanknote := func(tId string, banknote core.Money) {
		reqBody := api.BanknoteInsertionRequest{
			TransactionId: tId,
			Banknote:      api.Banknote{
				Currency:     string(banknote.Currency),
				Denomination: int(math.Round((banknote.Amount.Num()))),
			},
			Receivables:   []api.Money{}, // TODO: for now empty, maybe later fix to include receivables
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret)
		response := httptest.NewRecorder() 
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Deposit.OnBanknoteEscrow) 
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
	}
	deposit := func(user MockUser, tId string, dep core.Money) {
		reqBody := api.CompleteDepositRequest{
			TransactionId: tId,
			Receivables:   []api.Money{{
				string(dep.Currency),
				dep.Amount.Num(),
			}},
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret)
		response := httptest.NewRecorder()
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Deposit.OnComplete)
		handler.ServeHTTP(response, request) 
    AssertStatusCode(t, response, http.StatusOK)
	}
	newMoney := func(curr string, amount float64) core.Money {
		return core.Money{Currency: core.Currency(curr), Amount: core.NewMoneyAmount(amount)}
	}

	// assert that balance of user 1 is 0$
	assertBalance(users[0], newMoney("USD", 0))
	// assert that balance of user 2 is 0$
	assertBalance(users[1], newMoney("USD", 0))
	// deposit 100$ to user 1

	// assert that balance of user 1 is 100$
	// withdraw 49.5$ from user1
	// assert that balance of user1 is 50.5$
	// transfer 10.5$ from user1 to user2
	// assert that balance of user1 is 40$
	// assert that balance of user2 is 10.5$
	// withdraw 4.2$ from user2
	// assert that balance of user2 is 6.3$
	// deposit 5000 RUB to user2
	// assert that balance of user2 is 5000 RUB

}
