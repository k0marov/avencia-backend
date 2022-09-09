package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/di"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

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
		req.Header.Add("Authorization", authHeader)
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
	startTrans := func(tType values.TransactionType, qrText string) (tId string) {
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
	insertBanknote := func(tId string, b api.Banknote) {
		reqBody := api.BanknoteInsertionRequest{
			TransactionId: tId,
			Banknote:      b,
			Receivables:   []api.Money{}, // TODO: for now empty, maybe later fix to include receivables
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret)
		response := httptest.NewRecorder() 
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Deposit.OnBanknoteEscrow) 
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
	}
	finishDeposit := func(tId string, dep core.Money) {
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

	checkWithdrawal := func(tId string, b core.Money) {
		reqBody := api.StartWithdrawalRequest{
			TransactionId: tId,
			Currency:      string(b.Currency),
			Amount:        b.Amount.Num(),
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret) 
		response := httptest.NewRecorder()
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnStart)
		handler.ServeHTTP(response, request) 
    AssertStatusCode(t, response, http.StatusOK)
	}
	dispenseBanknote := func(tId string, b api.Banknote) {
		reqBody := api.BanknoteDispensionRequest{
			TransactionId:        tId,
			Currency:             b.Currency,
			BanknoteDenomination: b.Denomination,
			RemainingAmount:      float64(b.Denomination),
			RequestedAmount:      float64(b.Denomination),
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret)
		response := httptest.NewRecorder() 
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnPreBanknoteDispensed) 
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK) 
	}
	finishWithdrawal := func(tId string, w core.Money) {
		reqBody := api.CompleteWithdrawalRequest{
			TransactionId: tId,
			Currency:      string(w.Currency),
			Amount:        w.Amount.Num(),
		}
		request := addAuth(encodeRequest(reqBody), AtmAuthSecret) 
		response := httptest.NewRecorder() 
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnComplete)
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK)
	} 



	deposit := func(user MockUser, dep core.Money) {
		code := generateQRCode(user, values.Deposit) 
		// TODO: verify code.ExpiresAt
		tId := startTrans(values.Deposit, code.TransactionCode) 
		if dep.Amount.Num() > 1 {
			insertBanknote(tId, api.Banknote{
				Currency:     string(dep.Currency),
				Denomination: 1,
			}) 
		}
		finishDeposit(tId, dep) 
	}
	withdraw := func(user MockUser, w core.Money) {
	  code := generateQRCode(user, values.Withdrawal) 	
	  tId := startTrans(values.Withdrawal, code.TransactionCode)
	  checkWithdrawal(tId, w)
	  if w.Amount.Neg().Num() > 1 {
	  	dispenseBanknote(tId, api.Banknote{
	  		Currency:     string(w.Currency),
	  		Denomination: 1,
	  	}) 
	  }
	  finishWithdrawal(tId, w)
	}
	transfer := func(from, to MockUser, w core.Money) {
		reqBody := api.TransferRequest{
			RecipientIdentifier: to.Email,
			Money:               api.Money{
				Currency: string(w.Currency),
				Amount:   w.Amount.Num(),
			},
		}
		request := addAuth(encodeRequest(reqBody), from.Token) 
		response := httptest.NewRecorder() 
		handler := apiDeps.AtmAuthMW(apiDeps.Handlers.App.Transfer)
		handler.ServeHTTP(response, request) 
		AssertStatusCode(t, response, http.StatusOK)
	}

	// assert that balance of user 1 is 0$
	assertBalance(users[0], newMoney("USD", 0))
	// assert that balance of user 2 is 0$
	assertBalance(users[1], newMoney("USD", 0))
	// deposit 100$ to user 1
	deposit(users[0], newMoney("USD", 100))
	// assert that balance of user 1 is 100$
	assertBalance(users[0], newMoney("USD", 100))
	// withdraw 49.5$ from user1
	withdraw(users[0], newMoney("USD", 49.5))
	// assert that balance of user1 is 50.5$
	assertBalance(users[0], newMoney("USD", 50.5))
	// transfer 10.5$ from user1 to user2
	transfer(users[0], users[1], newMoney("USD", 10.5))
	// assert that balance of user1 is 40$
	assertBalance(users[0], newMoney("USD", 40.5))
	// assert that balance of user2 is 10.5$
	assertBalance(users[1], newMoney("USD", 10.5))
	// withdraw 4.2$ from user2
	withdraw(users[1], newMoney("USD", 4.2))
	// assert that balance of user2 is 6.3$
	assertBalance(users[1], newMoney("USD", 6.3))
	// deposit 5000 RUB to user2
	deposit(users[1], newMoney("RUB", 5000))
	// assert that balance of user2 is 5000 RUB
	assertBalance(users[1], newMoney("RUB", 5000))
}
