package http_helpers

import (
	"encoding/json"
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/service_helpers"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

type NoJSONRequest struct{}

type NoAPIResponse struct{}
func NoResponseConverter(service_helpers.Nothing) NoAPIResponse { return NoAPIResponse{} }

func NewAuthenticatedHandler[APIRequest any, Request any, Response any, APIResponse any](
	convertReq func(authEntities.User, *http.Request, APIRequest) (Request, error),
	service func(Request) (Response, error),
	convertResp func(Response) APIResponse,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		var req APIRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			ThrowHTTPError(w, client_errors.InvalidJSON)
			return
		}

		fullReq, err := convertReq(user, r, req)
		if err != nil {
			ThrowHTTPError(w, err)
			return
		}
		resp, err := service(fullReq)
		if err != nil {
			ThrowHTTPError(w, err)
			return
		}
		WriteJson(w, convertResp(resp))
	}
}

func NewHandler[APIRequest any, Request any, Response any, APIResponse any](
	convertReq func(*http.Request, APIRequest) (Request, error),
	service func(Request) (Response, error),
	convertResp func(Response) APIResponse,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req APIRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			ThrowHTTPError(w, client_errors.InvalidJSON)
			return
		}

		fullReq, err := convertReq(r, req)
		if err != nil {
			ThrowHTTPError(w, err)
			return
		}
		resp, err := service(fullReq)
		if err != nil {
			ThrowHTTPError(w, err)
			return
		}
		WriteJson(w, convertResp(resp))
	}
}
