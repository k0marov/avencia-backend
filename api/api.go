package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core"
	"net/http"
)

// TransactionTypeQueryArg Possible values: "deposit" and "withdrawal"
const TransactionTypeQueryArg = "transaction_type"
const TransactionTypeDeposit = "deposit"
const TransactionTypeWithdrawal = "withdrawal"

func NewAPIRouter(cashDeposit ATMTransaction, authMiddleware core.Middleware) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/atm-transaction", func(r chi.Router) {
			// requires a TransactionTypeQueryArg
			// Response: CodeResponse
			// Throws: TransactionTypeNotProvided
			r.Get("/gen-code", authMiddleware(cashDeposit.GenCode).ServeHTTP)

			// Request: CodeRequest; requires a TransactionTypeQueryArg
			// Response: VerifiedCodeResponse
			// Throws: TransactionTypeNotProvided, InvalidCode
			r.Post("/verify-code", cashDeposit.VerifyCode)

			// Request: BanknoteCheckRequest
			// Response: 200 if accepted, client error (or 500) if rejected
			// Throws: InvalidCode (means session-expired)
			r.Post("/check-banknote", cashDeposit.CheckBanknote)

			// Request: FinalizeTransactionRequest
			// Response: 200 if accepted, client error (or 500) if rejected
			// Throws: InvalidATMSecret
			r.Post("/finalize-transaction", cashDeposit.FinalizeTransaction)
		})
	})

	return r
}
