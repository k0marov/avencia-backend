package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core"
	"net/http"
)

type CashDepositHandlers struct {
	GenCode, VerifyCode, CheckBanknote, FinalizeTransaction http.HandlerFunc
}

func NewAPIRouter(cashDepositHandlers CashDepositHandlers, authMiddleware core.Middleware) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/cash/deposit", func(r chi.Router) {
			r.Get("/gen-code", authMiddleware(cashDepositHandlers.GenCode).ServeHTTP)
			r.Post("/verify-code", cashDepositHandlers.VerifyCode)
			r.Post("/check-banknote", cashDepositHandlers.CheckBanknote)
			r.Post("/finalize-transaction", cashDepositHandlers.FinalizeTransaction)
		})
	})

	return r
}
