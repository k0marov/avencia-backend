package deposit

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/router"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/service"
)

func NewDepositRouterImpl(authMiddleware core.Middleware) func(r chi.Router) {
	// service
	genCode := service.NewCodeGenerator(jwt.IssuerImpl)
	verifyCode := service.NewCodeVerifier(jwt.VerifierImpl)
	// handlers
	genCodeHandler := handlers.NewGenerateCodeHandler(genCode)
	verifyCodeHandler := handlers.NewVerifyCodeHandler(verifyCode)
	return router.NewDepositRouter(genCodeHandler, verifyCodeHandler, authMiddleware)
}
