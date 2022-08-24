package handlers

import (
	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	"net/http"
)

func NewTransferHandler(transfer service.DeliveryTransferer) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		apiRequests.TransferDecoder,
		http_helpers.NoResponseService(transfer),
		http_helpers.NoResponseConverter,
	)
}
