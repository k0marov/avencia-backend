package handlers

import (
	"encoding/json"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"log"
	"net/http"
)

func NewGenerateCodeHandler(generate service.CodeGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		transactionType := r.URL.Query().Get(api.TransactionTypeQueryArg)
		if transactionType == "" {
			http_helpers.ThrowClientError(w, client_errors.TransactionTypeNotProvided)
			return
		}
		code, err := generate(values.NewCode{
			TransType: values.TransactionType(transactionType),
			User:      user,
		})
		if err != nil {
			http_helpers.ThrowHTTPError(w, err)
			return
		}
		log.Printf("generated code %v for user %v", code, user.Id)
		http_helpers.WriteJson(w, api.CodeResponse{TransactionCode: code.Code, ExpiresAt: code.ExpiresAt.Unix()})
	}
}

func NewVerifyCodeHandler(verify service.CodeVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionType := r.URL.Query().Get(api.TransactionTypeQueryArg)
		if transactionType == "" {
			http_helpers.ThrowClientError(w, client_errors.TransactionTypeNotProvided)
			return
		}
		var code api.CodeRequest
		json.NewDecoder(r.Body).Decode(&code)
		userInfo, err := verify(values.CodeForCheck{
			Code:      code.TransactionCode,
			TransType: values.TransactionType(transactionType),
		})
		if err != nil {
			http_helpers.ThrowHTTPError(w, err)
			return
		}
		http_helpers.WriteJson(w, api.VerifiedCodeResponse{UserInfo: apiResponses.UserInfoEncoder(userInfo)})
	}
}

func NewCheckBanknoteHandler(checkBanknote service.BanknoteChecker) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.BanknoteDecoder,
		http_helpers.NoResponseService(checkBanknote),
		http_helpers.NoResponseConverter,
	)
}

func NewFinalizeTransactionHandler(finalizeTransaction service.ATMTransactionFinalizer) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.ATMTransactionDecoder,
		http_helpers.NoResponseService(finalizeTransaction),
		http_helpers.NoResponseConverter,
	)
}
