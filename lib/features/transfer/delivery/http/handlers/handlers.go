package handlers

import (
	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"net/http"
)

// TODO: maybe generalize the http handlers' structure, because they all look the same

func NewTransferHandler(transfer service.Transferer) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		apiRequests.TransferDecoder,
		http_helpers.NoResponseService(transfer),
		http_helpers.NoResponseConverter,
	)
}
