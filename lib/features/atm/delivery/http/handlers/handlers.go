package handlers

import (
	"net/http"

	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
)

func NewCreateTransactionHandler(create service.ATMTransactionCreator) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.TransDecoder,
		create,
		apiResponses.CreatedTransactionEncoder,
	)
}

func NewCancelTransactionHandler(cancel service.TransactionCanceler) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.CancelTransactionDecoder,
		http_helpers.NoResponseService(cancel),
		http_helpers.NoResponseConverter,
	)
}

func NewWithdrawalValidationHandler(validate validators.DeliveryWithdrawalValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
    apiRequests.WithdrawalDataDecoder, 
    http_helpers.NoResponseService(validate), 
    http_helpers.NoResponseConverter,
	)
}

func NewBanknoteEscrowHandler(validateBanknote validators.DeliveryInsertedBanknoteValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.InsertedBanknoteDecoder,
		http_helpers.NoResponseService(validateBanknote),
		http_helpers.NoResponseConverter,
	)
}

func NewBanknoteAcceptedHandler(validateBanknote validators.DeliveryInsertedBanknoteValidator) http.HandlerFunc {
	return NewBanknoteEscrowHandler(validateBanknote)
}

func NewPreBanknoteDispensedHandler(validateBanknote validators.DeliveryDispensedBanknoteValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.DispensedBanknoteDecoder,
		http_helpers.NoResponseService(validateBanknote),
		http_helpers.NoResponseConverter,
	)
}

func NewPostBanknoteDispensedHandler(validateBanknote validators.DeliveryDispensedBanknoteValidator) http.HandlerFunc {
	return NewPreBanknoteDispensedHandler(validateBanknote)
}
