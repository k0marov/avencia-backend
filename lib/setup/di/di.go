package di

import (
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/jwt"
	atmHandlers "github.com/AvenciaLab/avencia-backend/lib/features/atm/delivery/http/handlers"
	atmMiddleware "github.com/AvenciaLab/avencia-backend/lib/features/atm/delivery/http/middleware"
	atmService "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/service"
	atmValidators "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/validators"
	authMiddleware "github.com/AvenciaLab/avencia-backend/lib/features/auth/delivery/http/middleware"
	authService "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
	tStore "github.com/AvenciaLab/avencia-backend/lib/features/transactions/store"
	tService "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	histHandlers "github.com/AvenciaLab/avencia-backend/lib/features/histories/delivery/http/handlers"
	histEntities "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	histService "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
	histStore "github.com/AvenciaLab/avencia-backend/lib/features/histories/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	withdrawsService "github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/service"
	withdrawsStore "github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/store"
	tValidators "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/store/mappers"
	transHandlers "github.com/AvenciaLab/avencia-backend/lib/features/transfers/delivery/http/handlers"
	transService "github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/service"
	transValidators "github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/validators"
	userHandlers "github.com/AvenciaLab/avencia-backend/lib/features/users/delivery/http/handlers"
	userService "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	walletService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	storeImpl "github.com/AvenciaLab/avencia-backend/lib/features/wallets/store"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
	"github.com/go-chi/chi/v5/middleware"
)

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
	updBal := storeImpl.NewBalanceUpdater(db.JsonUpdaterImpl[core.MoneyAmount])
	getWallet := storeImpl.NewWalletGetter(db.JsonGetterImpl[walletEntities.Wallet])
	getBalance := walletService.NewBalanceGetter(getWallet)

	// ===== LIMITS =====
	storeGetWithdraws := withdrawsStore.NewWithdrawsGetter(db.JsonGetterImpl[models.Withdraws])
	storeUpdateWithdrawn := withdrawsStore.NewWithdrawUpdater(db.JsonUpdaterImpl[models.WithdrawVal])
	getUpdatedWithdrawn := withdrawsService.NewWithdrawnUpdateGetter(storeGetWithdraws)
	updateWithdrawn := withdrawsService.NewWithdrawnUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)
	getLimits := limits.NewLimitsGetter(storeGetWithdraws, limits.NewLimitsComputer(configurable.LimitedCurrencies))
	checkLimit := limits.NewLimitChecker(getLimits)

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
  
	transValidator := tValidators.NewTransactionValidator(checkLimit, tValidators.NewEnoughBalanceValidator(getBalance))
	codeParser := mappers.NewCodeParser(jwtVerifier)
	codeGenerator := mappers.NewCodeGenerator(jwtIssuer)

	createTrans := tStore.NewTransactionCreator(codeGenerator)
	getTrans := tStore.NewTransactionGetter(codeParser) 
	transPerformer := tService.NewTransactionPerformer(updateWithdrawn, storeTrans, tService.NewTransBalUpdater(updBal))
	transact := tService.NewTransactionFinalizer(transValidator, transPerformer)
	multiTransact := tService.NewMultiTransactionFinalizer(transact)

	// ===== ATM =====
	atmSecretValidator := atmValidators.NewATMSecretValidator(deps.AtmSecret)
	metaTransByIdValidator := atmValidators.NewMetaTransByIdValidator(getTrans)
	metaTransFromCodeValidator := atmValidators.NewMetaTransFromCodeValidator(codeParser)
	validateWithdrawal := atmValidators.NewWithdrawalValidator(metaTransByIdValidator, transValidator)
	insertedBanknoteValidator := atmValidators.NewInsertedBanknoteValidator()
	dispensedBanknoteValidator := atmValidators.NewDispensedBanknoteValidator()

	createAtmTrans := atmService.NewATMTransactionCreator(metaTransFromCodeValidator, createTrans)
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
