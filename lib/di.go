package lib

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade/batch"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	atmHandlers "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/store"
	"github.com/k0marov/avencia-backend/lib/features/transfers"
	userHandlers "github.com/k0marov/avencia-backend/lib/features/users/delivery/http/handlers"
	userService "github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/wallets/store"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"google.golang.org/api/option"
)

func initFirebase(config config.Config) *firebase.App {
	opt := option.WithCredentialsFile(config.FirebaseSecretPath)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	return fbApp
}

// TODO: write some integration tests (later)

func Initialize() http.Handler {
	conf := config.LoadConfig()
	// ===== CONFIG =====
	atmSecret, err := ioutil.ReadFile(conf.ATMSecretPath)
	if err != nil {
		log.Fatalf("error while reading atm secret: %v", err)
	}
	jwtSecret, err := ioutil.ReadFile(conf.JWTSecretPath)
	if err != nil {
		log.Fatalf("error while reading jwt secret: %v", err)
	}

	// ===== FIREBASE =====
	fbApp := initFirebase(conf)
	fsClient, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error while initializing firestore client: %v", err)
	}
	fbAuth, err := fbApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("erorr while initializing firebase auth: %v", err)
	}

	// ===== JWT =====
	jwtIssuer := jwt.NewIssuer(jwtSecret)
	jwtVerifier := jwt.NewVerifier(jwtSecret)

	// ===== AUTH =====
	authMiddleware := auth.NewAuthMiddleware(fbAuth)
	userFromEmail := auth.NewUserFromEmail(fbAuth)

	// ===== WALLET =====
	walletDocGetter := storeImpl.NewWalletDocGetter(firestore_facade.NewDocGetter(fsClient))
	storeGetWallet := storeImpl.NewWalletGetter(walletDocGetter)
	updateBalance := storeImpl.NewBalanceUpdater(walletDocGetter)
	getWallet := walletService.NewWalletGetter(storeGetWallet)
	getBalance := walletService.NewBalanceGetter(getWallet)

	// ===== LIMITS =====
	storeGetWithdraws := limitsStore.NewWithdrawsGetter(fsClient)
	storeUpdateWithdrawn := limitsStore.NewWithdrawUpdater(limitsStore.NewWithdrawDocGetter(firestore_facade.NewDocGetter(fsClient)))
	getLimits := limitsService.NewLimitsGetter(storeGetWithdraws, configurable.LimitedCurrencies)
	checkLimit := limitsService.NewLimitChecker(getLimits)
	getUpdatedWithdrawn := limitsService.NewWithdrawnUpdateGetter(getLimits)
	updateWithdrawn := limitsService.NewWithdrawUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)

	// ===== USER =====
	getUserInfo := userService.NewUserInfoGetter(getWallet, getLimits)
	getUserInfoHandler := userHandlers.NewGetUserInfoHandler(getUserInfo)

	// ===== ATM TRANSACTION =====
	// validators
	codeValidator := validators.NewTransCodeValidator(jwtVerifier)
	atmSecretValidator := validators.NewATMSecretValidator(atmSecret)
	transValidator := validators.NewTransactionValidator(checkLimit, getBalance)
	// service
	genCode := service.NewCodeGenerator(jwtIssuer)
	verifyCode := service.NewCodeVerifier(codeValidator, getUserInfo)
	checkBanknote := service.NewBanknoteChecker(verifyCode)
	performTrans := service.NewTransactionPerformer(updateBalance, updateWithdrawn)
	finalizeTransaction := service.NewTransactionFinalizer(transValidator, performTrans)
	atmFinalizeTransaction := service.NewATMTransactionFinalizer(atmSecretValidator, batch.NewWriteRunner(fsClient), finalizeTransaction)
	// handlers
	genCodeHandler := atmHandlers.NewGenerateCodeHandler(genCode)
	verifyCodeHandler := atmHandlers.NewVerifyCodeHandler(verifyCode)
	checkBanknoteHandler := atmHandlers.NewCheckBanknoteHandler(checkBanknote)
	atmTransactionHandler := atmHandlers.NewFinalizeTransactionHandler(atmFinalizeTransaction)

	// ===== TRANSFER =====
	transferHandler := transfers.NewTransferHandlerImpl(fsClient, userFromEmail, finalizeTransaction)

	apiRouter := api.NewAPIRouter(api.Handlers{
		GenCode:             genCodeHandler,
		VerifyCode:          verifyCodeHandler,
		CheckBanknote:       checkBanknoteHandler,
		FinalizeTransaction: atmTransactionHandler,
		GetUserInfo:         getUserInfoHandler,
		Transfer:            transferHandler,
	}, authMiddleware)
	return middleware.Recoverer(apiRouter)
}
