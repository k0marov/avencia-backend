package deposit

import (
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/deposit/delivery/http/router"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/service"
	"io/ioutil"
	"log"
)

func NewDepositRouterImpl(authMiddleware core.Middleware, config config.Config) func(r chi.Router) {
	// jwt
	jwtSecret, err := ioutil.ReadFile(config.JWTSecretPath)
	if err != nil {
		log.Fatalf("error while reading jwt secret: %v", err)
	}
	jwtIssuer := jwt.NewIssuer(jwtSecret)
	jwtVerifier := jwt.NewVerifier(jwtSecret)
	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	verifyCode := service.NewCodeVerifier(jwtVerifier)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	// handlers
	genCodeHandler := handlers.NewGenerateCodeHandler(genCode)
	verifyCodeHandler := handlers.NewVerifyCodeHandler(verifyCode)
	checkBanknoteHandler := handlers.NewCheckBanknoteHandler(checkBanknote)
	return router.NewDepositRouter(genCodeHandler, verifyCodeHandler, checkBanknoteHandler, authMiddleware)
}
