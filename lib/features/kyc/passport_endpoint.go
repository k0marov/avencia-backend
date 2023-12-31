package kyc

import (
	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-backend/lib/core/uploader"
	"github.com/go-chi/chi/v5"
)

func NewPassportEndpoint(upld uploader.UploaderFactory, statFactory StatusEndpointFactory) api.Endpoint {
	return func(r chi.Router) {
		r.Put("/front", upld("front", uploader.SimpleSizePolicy))
		r.Put("/back", upld("back", uploader.SimpleSizePolicy))
		r.Route("/status", statFactory("passport"))
	}
}
