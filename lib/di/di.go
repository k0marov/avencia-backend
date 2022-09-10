package di

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	atmHandlers "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/handlers"
	atmMiddleware "github.com/k0marov/avencia-backend/lib/features/atm/delivery/http/middleware"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	atmValidators "github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	authMiddleware "github.com/k0marov/avencia-backend/lib/features/auth/delivery/http/middleware"
	authService "github.com/k0marov/avencia-backend/lib/features/auth/domain/service"
	authStore "github.com/k0marov/avencia-backend/lib/features/auth/domain/store"
	histHandlers "github.com/k0marov/avencia-backend/lib/features/histories/delivery/http/handlers"
	histEntities "github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	histStore "github.com/k0marov/avencia-backend/lib/features/histories/store"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/store"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	transHandlers "github.com/k0marov/avencia-backend/lib/features/transfers/delivery/http/handlers"
	transService "github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	transValidators "github.com/k0marov/avencia-backend/lib/features/transfers/domain/validators"
	userHandlers "github.com/k0marov/avencia-backend/lib/features/users/delivery/http/handlers"
	userService "github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/wallets/store"
)

// TODO: write some integration tests (later)

type ExternalDeps struct {
	AtmSecret, JwtSecret []byte
	Auth                 authStore.AuthFacade
	TRunner              db.TransRunner
}
type APIDeps struct {
  Handlers api.Handlers 
  AuthMW api.Middleware
  AtmAuthMW api.Middleware
}


func InitializeBusiness(deps ExternalDeps) APIDeps {
	// ===== AUTH =====
	userAdder := authService.NewUserInfoAdder(deps.Auth.Verify)
	authMW := authMiddleware.NewAuthMiddleware(userAdder)
	userFromEmail := deps.Auth.UserByEmail

	// ===== JWT =====
	jwtIssuer := jwt.NewIssuer(deps.JwtSecret)
	jwtVerifier := jwt.NewVerifier(deps.JwtSecret)

	// ===== WALLETS =====
	updateBalance := storeImpl.NewBalanceUpdater(db.JsonUpdaterImpl[core.MoneyAmount])
	getWallet := storeImpl.NewWalletGetter(db.JsonGetterImpl[walletEntities.Wallet])
	getBalance := walletService.NewBalanceGetter(getWallet)

	// ===== LIMITS =====
	storeGetWithdraws := limitsStore.NewWithdrawsGetter(db.JsonGetterImpl[models.Withdraws])
	storeUpdateWithdrawn := limitsStore.NewWithdrawUpdater(db.JsonUpdaterImpl[models.WithdrawVal])
	getLimits := limitsService.NewLimitsGetter(storeGetWithdraws, configurable.LimitedCurrencies)
	checkLimit := limitsService.NewLimitChecker(getLimits)
	getUpdatedWithdrawn := limitsService.NewWithdrawnUpdateGetter(getLimits)
	updateWithdrawn := limitsService.NewWithdrawnUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)

	// ===== USERS =====
	getUserInfo := userService.NewUserInfoGetter(getWallet, getLimits)
	getUserInfoHandler := userHandlers.NewGetUserInfoHandler(deps.TRunner, getUserInfo)

	// ===== HISTORIES =====
	storeGetHistory := histStore.NewHistoryGetter(db.JsonColGetterImpl[histEntities.TransEntry])
	storeStoreTrans := histStore.NewTransStorer(db.JsonSetterImpl[histEntities.TransEntry])
	getHistory := histService.NewHistoryGetter(storeGetHistory)
	storeTrans := histService.NewTransStorer(storeStoreTrans)
	getHistoryHandler := histHandlers.NewGetHistoryHandler(deps.TRunner, getHistory)

	// ===== TRANSACTIONS =====
	transValidator := tValidators.NewTransactionValidator(checkLimit, getBalance)
	codeParser := mappers.NewCodeParser(jwtVerifier)
	codeGenerator := mappers.NewCodeGenerator(jwtIssuer)

	getTransId := tService.NewTransactionIdGetter(codeGenerator, mappers.NewTransIdGenerator())
	getTrans := tService.NewTransactionGetter(mappers.NewTransIdParser(), codeParser)
	transact := tService.NewTransactionFinalizer(transValidator, tService.NewTransactionPerformer(updateWithdrawn, storeTrans, updateBalance))
	multiTransact := tService.NewMultiTransactionFinalizer(transact)

	// TODO: write tests for the store layers

	// ===== ATM =====
	atmSecretValidator := atmValidators.NewATMSecretValidator(deps.AtmSecret)
	metaTransByIdValidator := atmValidators.NewMetaTransByIdValidator(getTrans)
	metaTransFromCodeValidator := atmValidators.NewMetaTransFromCodeValidator(codeParser)
	validateWithdrawal := atmValidators.NewWithdrawalValidator(metaTransByIdValidator, transValidator)
	insertedBanknoteValidator := atmValidators.NewInsertedBanknoteValidator()
	dispensedBanknoteValidator := atmValidators.NewDispensedBanknoteValidator()

	createAtmTrans := atmService.NewATMTransactionCreator(metaTransFromCodeValidator, getTransId)
	cancelTrans := atmService.NewTransactionCanceler()
	generalFinalizer := atmService.NewGeneralFinalizer(metaTransByIdValidator, multiTransact)
	finalizeDeposit := atmService.NewDepositFinalizer(generalFinalizer)
	finalizeWithdrawal := atmService.NewWithdrawalFinalizer(generalFinalizer)

	atmAuthMiddleware := atmMiddleware.NewATMAuthMiddleware(atmSecretValidator)

	genCodeHandler := atmHandlers.NewGenerateQRCodeHandler(codeGenerator)
	createTransHandler := atmHandlers.NewCreateTransactionHandler(createAtmTrans)
	onCancelHandler := atmHandlers.NewCancelTransactionHandler(cancelTrans)
	validateWithdrawalHandler := atmHandlers.NewWithdrawalValidationHandler(deps.TRunner, validateWithdrawal)
	completeDepositHandler := atmHandlers.NewCompleteDepostHandler(deps.TRunner, finalizeDeposit)
	completeWithdrawalHandler := atmHandlers.NewCompleteWithdrawalHandler(deps.TRunner, finalizeWithdrawal)
	banknoteEscrowHandler := atmHandlers.NewBanknoteEscrowHandler(deps.TRunner, insertedBanknoteValidator)
	banknoteAcceptedHandler := atmHandlers.NewBanknoteAcceptedHandler(deps.TRunner, insertedBanknoteValidator)
	preBanknoteDispensedHandler := atmHandlers.NewPreBanknoteDispensedHandler(deps.TRunner, dispensedBanknoteValidator)
	postBanknoteDispensedHandler := atmHandlers.NewPostBanknoteDispensedHandler(deps.TRunner, dispensedBanknoteValidator)

	// ===== TRANSFERS =====
	convertTransfer := transService.NewTransferConverter(userFromEmail)
	validateTransfer := transValidators.NewTransferValidator()
	performTransfer := transService.NewTransferPerformer(multiTransact)
	transfer := transService.NewTransferer(convertTransfer, validateTransfer, performTransfer)
	transferHandler := transHandlers.NewTransferHandler(deps.TRunner, transfer)

	return APIDeps{
		Handlers:  api.Handlers{
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
			GenCode:     genCodeHandler,
			GetUserInfo: getUserInfoHandler,
			Transfer:    transferHandler,
			GetHistory:  getHistoryHandler,
		},
	},
		AuthMW:    authMW,
		AtmAuthMW: atmAuthMiddleware,
	}
}

func InitializeHandler(deps APIDeps) http.Handler {
	apiRouter := api.NewAPIRouter(deps.Handlers, deps.AuthMW, deps.AtmAuthMW) 
	return middleware.Recoverer(apiRouter)
}
