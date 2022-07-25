package atm_transaction

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	userService "github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	"io/ioutil"
	"log"
	"net/http"
)

type WalletDeps struct {
	GetBalance    walletService.BalanceGetter
	UpdateBalance walletStore.BalanceUpdater
}

type UserDeps struct {
	GetUserInfo userService.UserInfoGetter
}

type LimitsDeps struct {
	CheckLimit          limitsService.LimitChecker
	GetUpdatedWithdrawn limitsService.WithdrawnUpdateGetter
	UpdateWithdrawn     limitsStore.WithdrawUpdater
}

type Handlers struct {
	GenCode, VerifyCode, CheckBanknote, FinalizeTransaction http.HandlerFunc
}

func NewATMTransactionHandlers(config config.Config, fsClient *firestore.Client, wallet WalletDeps, user UserDeps, limits LimitsDeps) Handlers {
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

	// validators
	codeValidator := validators.NewTransCodeValidator(jwtVerifier)
	transValidator := validators.NewTransactionValidator(atmSecret, limits.CheckLimit, wallet.GetBalance)

	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	verifyCode := service.NewCodeVerifier(codeValidator, user.GetUserInfo)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	performTrans := service.NewTransactionPerformer(batch.NewWriteRunner(fsClient), wallet.UpdateBalance, limits.GetUpdatedWithdrawn, limits.UpdateWithdrawn)
	finalizeTransaction := service.NewTransactionFinalizer(transValidator, performTrans)
	// handlers
	return Handlers{
		GenCode:             handlers.NewGenerateCodeHandler(genCode),
		VerifyCode:          handlers.NewVerifyCodeHandler(verifyCode),
		CheckBanknote:       handlers.NewCheckBanknoteHandler(checkBanknote),
		FinalizeTransaction: handlers.NewFinalizeTransactionHandler(finalizeTransaction),
	}
}
