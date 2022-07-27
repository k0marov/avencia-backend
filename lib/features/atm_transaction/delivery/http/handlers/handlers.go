package handlers

import (
	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"net/http"
)

func NewGenerateCodeHandler(generate service.CodeGenerator) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		apiRequests.NewCodeDecoder,
		generate,
		apiResponses.TransCodeEncoder,
	)
}

func NewVerifyCodeHandler(verify service.CodeVerifier) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.CodeForCheckDecoder,
		verify,
		apiResponses.UserInfoEncoder,
	)
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
