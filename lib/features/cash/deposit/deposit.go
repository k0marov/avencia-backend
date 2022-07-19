package deposit

import (
	"errors"
	"github.com/k0marov/avencia-backend/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/cash/deposit/domain/values"
	"io/ioutil"
	"log"
)

func NewCashDepositHandlers(config config.Config) api.CashDeposit {
	// jwt
	jwtSecret, err := ioutil.ReadFile(config.JWTSecretPath)
	if err != nil {
		log.Fatalf("error while reading jwt secret: %v", err)
	}
	jwtIssuer := jwt.NewIssuer(jwtSecret)
	jwtVerifier := jwt.NewVerifier(jwtSecret)

	// fake (not implemented yet)
	// TODO: implement transactionPerformer
	transactionPerformer := func(data values.TransactionData) error {
		log.Printf("fake performing a transaction: %+v", data)
		return errors.New("not implemented")
	}

	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	verifyCode := service.NewCodeVerifier(jwtVerifier)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	atmSecret, err := ioutil.ReadFile(config.ATMSecretPath)
	if err != nil {
		log.Fatalf("error while atm secret: %v", err)
	}
	finalizeTransaction := service.NewTransactionFinalizer(atmSecret, transactionPerformer)
	// handlers
	return api.CashDeposit{
		GenCode:             handlers.NewGenerateCodeHandler(genCode),
		VerifyCode:          handlers.NewVerifyCodeHandler(verifyCode),
		CheckBanknote:       handlers.NewCheckBanknoteHandler(checkBanknote),
		FinalizeTransaction: handlers.NewFinalizeTransactionHandler(finalizeTransaction),
	}
}
