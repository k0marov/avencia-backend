package handlers

import (
	"net/http"
	"net/url"

	"github.com/k0marov/avencia-api-contract/api"
	apiResponses "github.com/k0marov/avencia-backend/lib/api/api_responses"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func NewCreateTransactionHandler(create service.TransactionCreator) http.HandlerFunc {
  return http_helpers.NewHandler(
		func(_ url.Values, req api.OnTransactionCreateRequest) (values.NewTrans, error) {
			return values.NewTrans{
				Type:       tValues.TransactionType(req.Type),
				QRCodeText: req.QRCodeText,
			}, nil
		}, 
		create, 
		apiResponses.CreatedTransactionEncoder,
  )
}


