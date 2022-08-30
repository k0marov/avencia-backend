package handlers

import (
	"net/http"

	apiRequests "github.com/k0marov/avencia-backend/lib/api/api_requests"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/k0marov/avencia-backend/lib/core/helpers/service_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
)

func NewTransferHandler(runT db.TransRunner, transfer service.Transferer) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		apiRequests.TransferDecoder,
		service_helpers.NewDBNoResultService(runT, transfer), 
		http_helpers.NoResponseConverter,
	)
}
