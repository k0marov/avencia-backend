package atm_transaction

import (
	"errors"
	"github.com/k0marov/avencia-backend/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"io/ioutil"
	"log"
)

func NewATMTransactionHandlers(config config.Config) api.ATMTransaction {
	// config
	atmSecret, err := ioutil.ReadFile(config.ATMSecretPath)
	if err != nil {
		log.Fatalf("error while reading atm secret: %v", err)
	}

	// jwt
	jwtSecret, err := ioutil.ReadFile(config.JWTSecretPath)
	if err != nil {
		log.Fatalf("error while reading jwt secret: %v", err)
	}
	jwtIssuer := jwt.NewIssuer(jwtSecret)
	jwtVerifier := jwt.NewVerifier(jwtSecret)

	// fake (not implemented yet)
	// TODO: implement getBalance, updateBalance
	getBalance := func(userId string, currency string) (float64, error) {
		log.Printf("fake getting balance for user %s and currency %s", userId, currency)
		return 0, errors.New("unimplemented")
	}
	updateBalance := func(userId string, currency string, newBalance float64) error {
		log.Printf("fake setting balance for user %s and currency %s to %v", userId, currency, newBalance)
		return errors.New("unimplemented")
	}

	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	verifyCode := service.NewCodeVerifier(jwtVerifier)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	performTransaction := service.NewTransactionPerformer(getBalance, updateBalance)
	finalizeTransaction := service.NewTransactionFinalizer(atmSecret, performTransaction)
	// handlers
	return api.ATMTransaction{
		GenCode:             handlers.NewGenerateCodeHandler(genCode),
		VerifyCode:          handlers.NewVerifyCodeHandler(verifyCode),
		CheckBanknote:       handlers.NewCheckBanknoteHandler(checkBanknote),
		FinalizeTransaction: handlers.NewFinalizeTransactionHandler(finalizeTransaction),
	}
}
