package atm_transaction

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/store"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"io/ioutil"
	"log"
)

func NewATMTransactionHandlers(config config.Config, getWallet walletService.WalletGetter, fsClient *firestore.Client) api.ATMTransaction {
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

	// store
	getBalance := store.NewBalanceGetter(getWallet)
	performTrans := store.NewTransactionPerformer(fsClient)

	// validators
	codeValidator := validators.NewTransCodeValidator(jwtVerifier)
	transValidator := validators.NewTransactionValidator(atmSecret, getBalance)

	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	getUserInfo := service.NewUserInfoGetter(getWallet)
	verifyCode := service.NewCodeVerifier(codeValidator, getUserInfo)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	finalizeTransaction := service.NewTransactionFinalizer(transValidator, performTrans)
	// handlers
	return api.ATMTransaction{
		GenCode:             handlers.NewGenerateCodeHandler(genCode),
		VerifyCode:          handlers.NewVerifyCodeHandler(verifyCode),
		CheckBanknote:       handlers.NewCheckBanknoteHandler(checkBanknote),
		FinalizeTransaction: handlers.NewFinalizeTransactionHandler(finalizeTransaction),
	}
}
