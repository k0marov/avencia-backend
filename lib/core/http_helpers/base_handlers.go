package http_helpers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"net/http"
)

type NoResponse struct{}
type NoAPIResponse struct{}

func NoResponseConverter(NoResponse) NoAPIResponse { return NoAPIResponse{} }
func NoResponseService[Request any](service func(Request) error) func(Request) (NoResponse, error) {
	return func(request Request) (NoResponse, error) {
		return NoResponse{}, service(request)
	}
}

func NewAuthenticatedHandler[APIRequest any, Request any, Response any, APIResponse any](
	convertReq func(auth.User, APIRequest) Request,
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

		resp, err := service(convertReq(user, req))
		if err != nil {
			HandleServiceError(w, err)
			return
		}
		WriteJson(w, convertResp(resp))
	}
}

func NewHandler[APIRequest any, Request any, Response any, APIResponse any](
	convertReq func(APIRequest) Request,
	service func(Request) (Response, error),
	convertResp func(Response) APIResponse,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req APIRequest
		json.NewDecoder(r.Body).Decode(&req) // TODO: handle this error

		resp, err := service(convertReq(req))
		if err != nil {
			HandleServiceError(w, err)
			return
		}
		WriteJson(w, convertResp(resp))
	}
}
