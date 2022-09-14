package handlers

import (
	"net/http"

	apiRequests "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_requests"
	apiResponses "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_responses"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/mappers"
)

func NewGenerateQRCodeHandler(generate mappers.CodeGenerator) http.HandlerFunc {
  return http_helpers.NewAuthenticatedHandler(
  	apiRequests.NewTransDecoder, 
    generate, 
    apiResponses.TransCodeEncoder,
  ) 
}

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
		service_helpers.NewNoResultService(cancel),
		http_helpers.NoResponseConverter,
	)
}

func NewWithdrawalValidationHandler(runT db.TransRunner, validate validators.WithdrawalValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
    apiRequests.WithdrawalDataDecoder, 
    service_helpers.NewDBNoResultService(runT, validate),
    http_helpers.NoResponseConverter,
	)
}

func NewBanknoteEscrowHandler(runT db.TransRunner, validateBanknote validators.InsertedBanknoteValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.InsertedBanknoteDecoder,
		service_helpers.NewDBNoResultService(runT, validateBanknote),
		http_helpers.NoResponseConverter,
	)
}

func NewBanknoteAcceptedHandler(runT db.TransRunner, validateBanknote validators.InsertedBanknoteValidator) http.HandlerFunc {
	return NewBanknoteEscrowHandler(runT, validateBanknote)
}

// TODO: maybe use context.Context instead of db.DB for the services

func NewPreBanknoteDispensedHandler(runT db.TransRunner, validateBanknote validators.DispensedBanknoteValidator) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.DispensedBanknoteDecoder,
		service_helpers.NewDBNoResultService(runT, validateBanknote),
		http_helpers.NoResponseConverter,
	)
}

func NewPostBanknoteDispensedHandler(runT db.TransRunner, validateBanknote validators.DispensedBanknoteValidator) http.HandlerFunc {
	return NewPreBanknoteDispensedHandler(runT, validateBanknote)
}


func NewCompleteDepostHandler(runT db.TransRunner, finalize service.DepositFinalizer) http.HandlerFunc {
	return http_helpers.NewHandler(
    apiRequests.DepositDataDecoder, 
    service_helpers.NewDBNoResultService(runT, finalize),   
    http_helpers.NoResponseConverter, 
	)
}

func NewCompleteWithdrawalHandler(runT db.TransRunner, finalize service.WithdrawalFinalizer) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.WithdrawalDataDecoder, 
    service_helpers.NewDBNoResultService(runT, finalize),   
		http_helpers.NoResponseConverter, 
	)
}


