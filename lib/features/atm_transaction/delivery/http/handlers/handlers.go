package handlers

import (
	"encoding/json"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
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
		code, expiresAt, err := generate(user, values.TransactionType(transactionType))
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		log.Printf("generated code %v for user %v", code, user.Id)
		http_helpers.WriteJson(w, api.CodeResponse{TransactionCode: code, ExpiresAt: expiresAt.Unix()})
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
		userInfo, err := verify(code.TransactionCode, values.TransactionType(transactionType))
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, api.VerifiedCodeResponse{UserInfo: userInfo.ToResponse()})
	}
}

func NewCheckBanknoteHandler(checkBanknote service.BanknoteChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var banknoteRequest api.BanknoteCheckRequest
		json.NewDecoder(r.Body).Decode(&banknoteRequest)

		err := checkBanknote(banknoteRequest.TransactionCode, values.NewBanknote(banknoteRequest))
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
	}
}

func NewFinalizeTransactionHandler(finalizeTransaction service.ATMTransactionFinalizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t api.FinalizeTransactionRequest
		json.NewDecoder(r.Body).Decode(&t)

		err := finalizeTransaction([]byte(t.ATMSecret), values.NewTransactionData(t))
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
	}
}
