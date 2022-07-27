package http_helpers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"net/http"
	"net/url"
)

type NoJSONRequest struct{}

type NoResponse struct{}
type NoAPIResponse struct{}

func NoResponseConverter(NoResponse) NoAPIResponse { return NoAPIResponse{} }
func NoResponseService[Request any](service func(Request) error) func(Request) (NoResponse, error) {
	return func(request Request) (NoResponse, error) {
		return NoResponse{}, service(request)
	}
}

func NewAuthenticatedHandler[APIRequest any, Request any, Response any, APIResponse any](
	convertReq func(auth.User, url.Values, APIRequest) (Request, error),
	service func(Request) (Response, error),
	convertResp func(Response) APIResponse,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		var req APIRequest
		json.NewDecoder(r.Body).Decode(&req) // TODO: handle this error

		fullReq, err := convertReq(user, r.URL.Query(), req)
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
	convertReq func(url.Values, APIRequest) (Request, error),
	service func(Request) (Response, error),
	convertResp func(Response) APIResponse,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req APIRequest
		json.NewDecoder(r.Body).Decode(&req) // TODO: handle this error
		fullReq, err := convertReq(r.URL.Query(), req)
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
