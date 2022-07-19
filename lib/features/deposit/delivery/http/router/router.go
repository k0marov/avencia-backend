package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core"
	"net/http"
)

func NewDepositRouter(generateCode, verifyCode, checkBanknote, finalizeTransaction http.HandlerFunc, authMiddleware core.Middleware) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/gen-code", authMiddleware(generateCode).ServeHTTP)
		r.Post("/verify-code", verifyCode)
		r.Post("/check-banknote", checkBanknote)
		r.Post("/finalize-transaction", finalizeTransaction)
	}
}
