package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core"
	"net/http"
)

func NewDepositRouter(generateCode, verifyCode http.HandlerFunc, authMiddleware core.Middleware) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/gen-code", authMiddleware(generateCode).ServeHTTP)
		r.Get("/verify-code", verifyCode)
	}
}
