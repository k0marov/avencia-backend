package handlers

import (
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/currencies/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/setup/api/api_requests"
	"github.com/AvenciaLab/avencia-backend/lib/setup/api/api_responses"
)

func NewGetExchangeRatesHandler(getRates service.ExchangeRatesGetter) http.HandlerFunc {
	return http_helpers.NewHandler(
		apiRequests.CurrenciesDecoder,
		getRates,
		apiResponses.ExchangeRatesEncoder,
	)
}
