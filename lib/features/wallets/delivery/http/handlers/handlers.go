package handlers

import (
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	apiRequests "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_requests"
	apiResponses "github.com/AvenciaLab/avencia-backend/lib/setup/api/api_responses"
)

func NewGetWalletsHandler(runT db.TransRunner, getWallets service.WalletsGetter) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		func(user authEntities.User, _ *http.Request, _ http_helpers.NoJSONRequest) (string, error) { return user.Id, nil },
		service_helpers.NewDBTransService(runT, getWallets),
		apiResponses.WalletsEncoder,
	)
}

func NewCreateWalletHandler(runT db.TransRunner, create service.WalletCreator) http.HandlerFunc {
	return http_helpers.NewAuthenticatedHandler(
		apiRequests.WalletCreationDecoder,
		service_helpers.NewDBTransService(runT, create),
		apiResponses.IdEncoder,
	)
}
