package lib

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	atmHandlers "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/handlers"
	atmMiddleware "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/middleware"
	atmValidators "github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	histHandlers "github.com/k0marov/avencia-backend/lib/features/histories/delivery/http/handlers"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	histStore "github.com/k0marov/avencia-backend/lib/features/histories/store"
	histMappers "github.com/k0marov/avencia-backend/lib/features/histories/store/mappers"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/store"
	limitsMappers "github.com/k0marov/avencia-backend/lib/features/limits/store/mappers"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	transHandlers "github.com/k0marov/avencia-backend/lib/features/transfers/delivery/http/handlers"
	transService "github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	transValidators "github.com/k0marov/avencia-backend/lib/features/transfers/domain/validators"
	userHandlers "github.com/k0marov/avencia-backend/lib/features/users/delivery/http/handlers"
	userService "github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/wallets/store"

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
	// jwtSecret, err := ioutil.ReadFile(conf.JWTSecretPath)
	// if err != nil {
	// 	log.Fatalf("error while reading jwt secret: %v", err)
	// }

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

	// ===== DB =====

	// ===== JWT =====
	// jwtIssuer := jwt.NewIssuer(jwtSecret)
	// jwtVerifier := jwt.NewVerifier(jwtSecret)

	// ===== AUTH =====
	authMiddleware := auth.NewAuthMiddleware(fbAuth)
	userFromEmail := auth.NewUserFromEmail(fbAuth)

	// ===== WALLETS =====
	walletDocGetter := storeImpl.NewWalletDocGetter(fsDocGetter)
	storeGetWallet := storeImpl.NewWalletGetter(walletDocGetter)
	updateBalance := storeImpl.NewBalanceUpdater(walletDocGetter)
	getWallet := walletService.NewWalletGetter(storeGetWallet)
	getBalance := walletService.NewBalanceGetter(getWallet)

	// ===== LIMITS =====
	storeGetWithdraws := limitsStore.NewWithdrawsGetter(fsClient, limitsMappers.WithdrawsDecoderImpl)
	storeUpdateWithdrawn := limitsStore.NewWithdrawUpdater(fsDocGetter, limitsMappers.WithdrawEncoderImpl)
	getLimits := limitsService.NewLimitsGetter(storeGetWithdraws, configurable.LimitedCurrencies)
	checkLimit := limitsService.NewLimitChecker(getLimits)
	getUpdatedWithdrawn := limitsService.NewWithdrawnUpdateGetter(getLimits)
	updateWithdrawn := limitsService.NewWithdrawnUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)

	// ===== USERS =====
	getUserInfo := userService.NewUserInfoGetter(getWallet, getLimits)
	getUserInfoHandler := userHandlers.NewGetUserInfoHandler(getUserInfo)

	// ===== HISTORIES =====
	storeGetHistory := histStore.NewHistoryGetter(fsClient, histMappers.TransEntriesDecoderImpl)
	storeStoreTrans := histStore.NewTransStorer(fsDocGetter, histMappers.TransEntryEncoderImpl)
	getHistory := histService.NewHistoryGetter(storeGetHistory)
	storeTrans := histService.NewTransStorer(storeStoreTrans)
	getHistoryHandler := histHandlers.NewGetHistoryHandler(getHistory)

	// ===== TRANSACTIONS =====
	transValidator := tValidators.NewTransactionValidator(checkLimit, getBalance)
	transact := tService.NewTransactionFinalizer(transValidator, tService.NewTransactionPerformer(updateWithdrawn, storeTrans, updateBalance))

	// ===== ATM =====
	atmSecretValidator := atmValidators.NewATMSecretValidator(atmSecret)
	// codeValidator := atmValidators.NewTransCodeValidator(jwtVerifier)
	
	atmAuthMiddleware := atmMiddleware.NewATMAuthMiddleware(atmSecretValidator)

	createTransHandler := atmHandlers.NewCreateTransactionHandler(nil) // TODO: add service 



	// ===== TRANSFERS =====
	convertTransfer := transService.NewTransferConverter(userFromEmail)
	validateTransfer := transValidators.NewTransferValidator()
	performTransfer := transService.NewTransferPerformer(runBatch, transact)
	transfer := transService.NewTransferer(convertTransfer, validateTransfer, performTransfer)
	transferHandler := transHandlers.NewTransferHandler(transfer)

	apiRouter := api.NewAPIRouter(api.Handlers{
		Transaction: api.TransactionHandlers{
			OnCreate: createTransHandler,
			OnCancel: nil,
			Deposit: api.TransactionDepositHandlers{
				OnBanknoteEscrow:   nil,
				OnBanknoteAccepted: nil,
				OnComplete:         nil,
			},
			Withdrawal: api.TransactionWithdrawalHandlers{
				OnStart:                 nil,
				OnPreBanknoteDispensed:  nil,
				OnPostBanknoteDispensed: nil,
				OnComplete:              nil,
			},
		},
		App: api.AppHandlers{
			GenCode:     nil,
			GetUserInfo: getUserInfoHandler,
			Transfer:    transferHandler,
			GetHistory:  getHistoryHandler,
		},
	}, authMiddleware, atmAuthMiddleware)
	return middleware.Recoverer(apiRouter)
}
