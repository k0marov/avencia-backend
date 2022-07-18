package handlers

import (
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/responses"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/service"
	"net/http"
)

func NewGenerateCodeHandler(generate service.CodeGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		code, err := generate(user)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, responses.CodeResponse{Code: code})
	}
}

type CodeRequest struct {
	Code string `json:"code"`
}

func NewVerifyCodeHandler(verify service.CodeVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code CodeRequest
		json.NewDecoder(r.Body).Decode(&code)
		userInfo, err := verify(code.Code)
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
		http_helpers.WriteJson(w, responses.UserInfoResponse{Id: userInfo.Id})
	}
}
