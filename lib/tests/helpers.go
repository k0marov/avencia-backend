package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

// apiDeps is for simplicity a global variable instead of an argument to every helper
var (
	apiDeps di.APIDeps 
	atmAuthSecret string
  jwtSecret string 
)


func initApiDeps(deps di.APIDeps, atmAuthSecr, jwtSecr string) {
	apiDeps = deps 
	atmAuthSecret = atmAuthSecr
	jwtSecret = jwtSecr
}

func deposit(t *testing.T, user MockUser, dep core.Money) {
	code := generateQRCode(t, user, values.Deposit)
	// TODO: verify code.ExpiresAt
	tId := startTrans(t, values.Deposit, code.TransactionCode)
	if dep.Amount.Num() > 1 {
		insertBanknote(t, tId, api.Banknote{
			Currency:     string(dep.Currency),
			Denomination: 1,
		})
	}
	finishDeposit(t, tId, dep)
}
func withdraw(t *testing.T, user MockUser, w core.Money) {
	code := generateQRCode(t, user, values.Withdrawal)
	tId := startTrans(t, values.Withdrawal, code.TransactionCode)
	checkWithdrawal(t, tId, w)
	if w.Amount.Neg().Num() > 1 {
		dispenseBanknote(t, tId, api.Banknote{
			Currency:     string(w.Currency),
			Denomination: 1,
		})
	}
	finishWithdrawal(t, tId, w)
}
func transfer(t *testing.T, from, to MockUser, w core.Money) {
	reqBody := api.TransferRequest{
		RecipientIdentifier: to.Email,
		Money: api.Money{
			Currency: string(w.Currency),
			Amount:   w.Amount.Num(),
		},
	}
	request := addAuth(encodeRequest(t, reqBody), from.Token)
	response := httptest.NewRecorder()
	handler := apiDeps.AuthMW(apiDeps.Handlers.App.Transfer)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}
func assertBalance(t *testing.T, user MockUser, balance core.Money) {
	request := addAuth(encodeRequest(t, nil), user.Token)
	response := httptest.NewRecorder()
	handler := apiDeps.AuthMW(apiDeps.Handlers.App.GetUserInfo)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
	var userInfo api.UserInfoResponse
	err := json.Unmarshal(response.Body.Bytes(), &userInfo)
	AssertNoError(t, err)
	Assert(t, userInfo.Wallet[string(balance.Currency)], balance.Amount.Num(), "balance")
}


func encodeRequest(t testing.TB, req any) *http.Request {
	t.Helper()
	body, err := json.Marshal(req)
	AssertNoError(t, err)
	return httptest.NewRequest("", "/asdf", bytes.NewReader(body))
}
func addAuth(req *http.Request, authHeader string) *http.Request {
	req.Header.Add("Authorization", authHeader)
	return req
}

func generateQRCode(t testing.TB, user MockUser, transType values.TransactionType) api.GenTransCodeResponse {
	t.Helper()
	reqBody := api.GenTransCodeRequest{
		TransactionType: string(transType),
	}
	request := addAuth(encodeRequest(t, reqBody), user.Token)
	response := httptest.NewRecorder()
	handler := apiDeps.AuthMW(apiDeps.Handlers.App.GenCode)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
	var codeResp api.GenTransCodeResponse
	err := json.Unmarshal(response.Body.Bytes(), &codeResp)
	AssertNoError(t, err)
	return codeResp
}
func startTrans(t testing.TB, tType values.TransactionType, qrText string) (tId string) {
	t.Helper()
	reqBody := api.OnTransactionCreateRequest{
		TransactionReference: "asdf",
		TerminalId:           "1234",
		TerminalSid:          "4321",
		Type:                 string(tType),
		QRCodeText:           qrText,
	}
	response := httptest.NewRecorder()
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.OnCreate)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
	var jsonResp api.OnTransactionCreateResponse
	err := json.Unmarshal(response.Body.Bytes(), &jsonResp)
	AssertNoError(t, err)
	return jsonResp.TransactionId
}
func insertBanknote(t testing.TB, tId string, b api.Banknote) {
	t.Helper()
	reqBody := api.BanknoteInsertionRequest{
		TransactionId: tId,
		Banknote:      b,
		Receivables:   []api.Money{}, // TODO: for now empty, maybe later fix to include receivables
	}
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	response := httptest.NewRecorder()
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Deposit.OnBanknoteEscrow)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}
func finishDeposit(t testing.TB, tId string, dep core.Money) {
	t.Helper()
	reqBody := api.CompleteDepositRequest{
		TransactionId: tId,
		Receivables: []api.Money{{
			Currency: string(dep.Currency),
			Amount: dep.Amount.Num(),
		}},
	}
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	response := httptest.NewRecorder()
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Deposit.OnComplete)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}
func newMoney(curr string, amount float64) core.Money {
	return core.Money{Currency: core.Currency(curr), Amount: core.NewMoneyAmount(amount)}
}

func checkWithdrawal(t testing.TB, tId string, b core.Money) {
	t.Helper()
	reqBody := api.StartWithdrawalRequest{
		TransactionId: tId,
		Currency:      string(b.Currency),
		Amount:        b.Amount.Num(),
	}
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	response := httptest.NewRecorder()
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnStart)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}
func dispenseBanknote(t testing.TB, tId string, b api.Banknote) {
	t.Helper()
	reqBody := api.BanknoteDispensionRequest{
		TransactionId:        tId,
		Currency:             b.Currency,
		BanknoteDenomination: b.Denomination,
		RemainingAmount:      float64(b.Denomination),
		RequestedAmount:      float64(b.Denomination),
	}
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	response := httptest.NewRecorder()
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnPreBanknoteDispensed)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}
func finishWithdrawal(t testing.TB, tId string, w core.Money) {
	t.Helper()
	reqBody := api.CompleteWithdrawalRequest{
		TransactionId: tId,
		Currency:      string(w.Currency),
		Amount:        w.Amount.Num(),
	}
	request := addAuth(encodeRequest(t, reqBody), atmAuthSecret)
	response := httptest.NewRecorder()
	handler := apiDeps.AtmAuthMW(apiDeps.Handlers.Transaction.Withdrawal.OnComplete)
	handler.ServeHTTP(response, request)
	AssertStatusCode(t, response, http.StatusOK)
}

