package handlers

import (
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
		}
		http_helpers.WriteJson(w, responses.CodeResponse{Code: code})
	}
}

func NewVerifyCodeHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {

	}
}
