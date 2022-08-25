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
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/db/firestore"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	atmHandlers "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/handlers"
	atmMiddleware "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/middleware"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	atmValidators "github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	histHandlers "github.com/k0marov/avencia-backend/lib/features/histories/delivery/http/handlers"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	histStore "github.com/k0marov/avencia-backend/lib/features/histories/store"
	histMappers "github.com/k0marov/avencia-backend/lib/features/histories/store/mappers"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/store"
	limitsMappers "github.com/k0marov/avencia-backend/lib/features/limits/store/mappers"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
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
		log.Fatalf("error while initializing firebase auth: %v", err)
	}

	// ===== DB =====
	simpleDB := db.NewDB(firestore.NewSimpleDB(fsClient))
	runTransaction := firestore.NewTransactionRunner(fsClient)

	// ===== JWT =====
	jwtIssuer := jwt.NewIssuer(jwtSecret)
	jwtVerifier := jwt.NewVerifier(jwtSecret)

	// ===== AUTH =====
	authMiddleware := auth.NewAuthMiddleware(fbAuth)
	userFromEmail := auth.NewUserFromEmail(fbAuth)

	// ===== WALLETS =====
	storeGetWallet := storeImpl.NewWalletGetter(db.DocGetterImpl)
	updateBalance := storeImpl.NewBalanceUpdater(db.SetterImpl)
	getWallet := walletService.NewWalletGetter(storeGetWallet)
	getBalance := walletService.NewBalanceGetter(getWallet)

	// ===== LIMITS =====
	storeGetWithdraws := limitsStore.NewWithdrawsGetter(db.ColGetterImpl, limitsMappers.WithdrawsDecoderImpl)
	storeUpdateWithdrawn := limitsStore.NewWithdrawUpdater(db.SetterImpl, limitsMappers.WithdrawEncoderImpl)
	getLimits := limitsService.NewLimitsGetter(storeGetWithdraws, configurable.LimitedCurrencies)
	checkLimit := limitsService.NewLimitChecker(getLimits)
	getUpdatedWithdrawn := limitsService.NewWithdrawnUpdateGetter(getLimits)
	updateWithdrawn := limitsService.NewWithdrawnUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)

	// ===== USERS =====
	getUserInfo := userService.NewUserInfoGetter(getWallet, getLimits)
	getUserInfoDelivery := userService.NewDeliveryUserInfoGetter(simpleDB, getUserInfo)
	getUserInfoHandler := userHandlers.NewGetUserInfoHandler(getUserInfoDelivery)

	// ===== HISTORIES =====
	storeGetHistory := histStore.NewHistoryGetter(db.ColGetterImpl, histMappers.TransEntriesDecoderImpl)
	storeStoreTrans := histStore.NewTransStorer(db.SetterImpl, histMappers.TransEntryEncoderImpl)
	getHistory := histService.NewHistoryGetter(storeGetHistory)
	getHistoryDelivery := histService.NewDeliveryHistoryGetter(simpleDB, getHistory)
	storeTrans := histService.NewTransStorer(storeStoreTrans)
	getHistoryHandler := histHandlers.NewGetHistoryHandler(getHistoryDelivery)

	// ===== TRANSACTIONS =====
	transValidator := tValidators.NewTransactionValidator(checkLimit, getBalance)
	codeParser := mappers.NewCodeParser(jwtVerifier)
	codeGenerator := mappers.NewCodeGenerator(jwtIssuer)

	getTransId := tService.NewTransactionIdGetter(codeGenerator, mappers.NewTransIdGenerator())
	getTrans := tService.NewTransactionGetter(mappers.NewTransIdParser(), codeParser)
	transact := tService.NewTransactionFinalizer(transValidator, tService.NewTransactionPerformer(updateWithdrawn, storeTrans, updateBalance))
	multiTransact := tService.NewMultiTransactionFinalizer(transact) 

	// ===== ATM =====
	atmSecretValidator := atmValidators.NewATMSecretValidator(atmSecret)
	metaTransByIdValidator := atmValidators.NewMetaTransByIdValidator(getTrans) 
	metaTransFromCodeValidator := atmValidators.NewMetaTransFromCodeValidator(codeParser)
	validateWithdrawal := atmService.NewDeliveryWithdrawalValidator(simpleDB, atmValidators.NewWithdrawalValidator(metaTransByIdValidator, transValidator))
	insertedBanknoteValidator := atmService.NewDeliveryInsertedBanknoteValidator(simpleDB, atmValidators.NewInsertedBanknoteValidator())
	dispensedBanknoteValidator := atmService.NewDeliveryDispensedBanknoteValidator(simpleDB, atmValidators.NewDispensedBanknoteValidator())

	createAtmTrans  := atmService.NewATMTransactionCreator(metaTransFromCodeValidator, getTransId)
	cancelTrans := atmService.NewTransactionCanceler()
	generalFinalizer := atmService.NewGeneralFinalizer(metaTransByIdValidator, multiTransact)
	finalizeDeposit := atmService.NewDeliveryDepositFinalizer(runTransaction, atmService.NewDepositFinalizer(generalFinalizer))
	finalizeWithdrawal := atmService.NewDeliveryWithdrawalFinalizer(runTransaction, atmService.NewWithdrawalFinalizer(generalFinalizer))
	
	
	atmAuthMiddleware := atmMiddleware.NewATMAuthMiddleware(atmSecretValidator)

	createTransHandler := atmHandlers.NewCreateTransactionHandler(createAtmTrans) 
	onCancelHandler := atmHandlers.NewCancelTransactionHandler(cancelTrans)
	validateWithdrawalHandler := atmHandlers.NewWithdrawalValidationHandler(validateWithdrawal)
	completeDepositHandler := atmHandlers.NewCompleteDepostHandler(finalizeDeposit) 
	completeWithdrawalHandler := atmHandlers.NewCompleteWithdrawalHandler(finalizeWithdrawal)
	banknoteEscrowHandler := atmHandlers.NewBanknoteEscrowHandler(insertedBanknoteValidator)
	banknoteAcceptedHandler := atmHandlers.NewBanknoteAcceptedHandler(insertedBanknoteValidator)
	preBanknoteDispensedHandler := atmHandlers.NewPreBanknoteDispensedHandler(dispensedBanknoteValidator)
	postBanknoteDispensedHandler := atmHandlers.NewPostBanknoteDispensedHandler(dispensedBanknoteValidator)
	


	// ===== TRANSFERS =====
	convertTransfer := transService.NewTransferConverter(userFromEmail)
	validateTransfer := transValidators.NewTransferValidator()
	performTransfer := transService.NewTransferPerformer(multiTransact)
	transfer := transService.NewTransferer(convertTransfer, validateTransfer, performTransfer)
	transferDelivery := transService.NewDeliveryTransferer(runTransaction, transfer)
	transferHandler := transHandlers.NewTransferHandler(transferDelivery)

	apiRouter := api.NewAPIRouter(api.Handlers{
		Transaction: api.TransactionHandlers{
			OnCreate: createTransHandler,
			OnCancel: onCancelHandler,
			Deposit: api.TransactionDepositHandlers{
				OnBanknoteEscrow:   banknoteEscrowHandler,
				OnBanknoteAccepted: banknoteAcceptedHandler,
				OnComplete:         completeDepositHandler,
			},
			Withdrawal: api.TransactionWithdrawalHandlers{
				OnStart:                 validateWithdrawalHandler,
				OnPreBanknoteDispensed:  preBanknoteDispensedHandler,
				OnPostBanknoteDispensed: postBanknoteDispensedHandler,
				OnComplete:              completeWithdrawalHandler,
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
