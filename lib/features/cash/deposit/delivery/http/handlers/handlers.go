package handlers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/api"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/values"
	"log"
	"net/http"
)

func NewGenerateCodeHandler(generate service.CodeGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		code, expiresAt, err := generate(user)
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
		var code api.CodeRequest
		json.NewDecoder(r.Body).Decode(&code)
		userInfo, err := verify(code.TransactionCode)
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

		accept := checkBanknote(banknoteRequest.TransactionCode, values.NewBanknote(banknoteRequest))
		response := api.AcceptionResponse{Accept: accept}

		http_helpers.WriteJson(w, response)
	}
}

func NewFinalizeTransactionHandler(finalizeTransaction service.TransactionFinalizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transactionRequest api.FinalizeTransactionRequest
		json.NewDecoder(r.Body).Decode(&transactionRequest)

		accept := finalizeTransaction(values.NewTransactionData(transactionRequest))
		response := api.AcceptionResponse{Accept: accept}

		http_helpers.WriteJson(w, response)
	}
}
