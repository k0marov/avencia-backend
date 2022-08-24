package handlers

import (
	"net/http"

	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
)

func NewCreateTransactionHandler(create service.ATMTransactionCreator) http.HandlerFunc {
  return http_helpers.NewHandler(
  	apiRequests.NewTransDecoder,
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


