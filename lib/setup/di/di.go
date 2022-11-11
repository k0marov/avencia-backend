package di

import (
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/jwt"
	"github.com/AvenciaLab/avencia-backend/lib/core/static_store"
	"github.com/AvenciaLab/avencia-backend/lib/core/uploader"
	atmHandlers "github.com/AvenciaLab/avencia-backend/lib/features/atm/delivery/http/handlers"
	atmMiddleware "github.com/AvenciaLab/avencia-backend/lib/features/atm/delivery/http/middleware"
	atmService "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/service"
	atmValidators "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/validators"
	authMiddleware "github.com/AvenciaLab/avencia-backend/lib/features/auth/delivery/http/middleware"
	authService "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	histHandlers "github.com/AvenciaLab/avencia-backend/lib/features/histories/delivery/http/handlers"
	histEntities "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	histService "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
	histStore "github.com/AvenciaLab/avencia-backend/lib/features/histories/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/kyc"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/models"
	withdrawsService "github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/service"
	withdrawsStore "github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	tService "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	tStore "github.com/AvenciaLab/avencia-backend/lib/features/transactions/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/store/mappers"
	"github.com/AvenciaLab/avencia-backend/lib/features/users"
	uHandlers "github.com/AvenciaLab/avencia-backend/lib/features/users/delivery/http/handlers"
	userService "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
	wHandlers "github.com/AvenciaLab/avencia-backend/lib/features/wallets/delivery/http/handlers"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	walletService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	storeImpl "github.com/AvenciaLab/avencia-backend/lib/features/wallets/store"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
	"github.com/go-chi/chi/v5/middleware"
)

type ExternalDeps struct {
	Config               config.Config
	AtmSecret, JwtSecret []byte
	Auth                 authStore.AuthFacade
	TRunner              db.TransRunner
	SimpleDB             db.SDB
}
type APIDeps struct {
	Handlers  api.Handlers
	AuthMW    api.Middleware
	AtmAuthMW api.Middleware
}

func InitializeBusiness(deps ExternalDeps) APIDeps {
	// ===== AUTH =====
	userAdder := authService.NewUserInfoAdder(deps.Auth.Verify)
	authMW := authMiddleware.NewAuthMiddleware(userAdder)
	// userFromEmail := deps.Auth.UserByEmail

	// ===== JWT =====
	jwtIssuer := jwt.NewIssuer(deps.JwtSecret)
	jwtVerifier := jwt.NewVerifier(deps.JwtSecret)

	// ===== WALLETS =====
	updBal := walletService.NewBalanceUpdater(storeImpl.NewBalanceUpdater(db.JsonUpdaterImpl[core.MoneyAmount]))
	getWallet := walletService.NewWalletGetter(storeImpl.NewWalletGetter(db.JsonGetterImpl[walletEntities.WalletVal]))
	storeCreateWallet := storeImpl.NewWalletCreator(
		db.JsonGetterImpl[storeImpl.UserWalletsModel], 
		db.JsonUpdaterImpl[[]string], 
    db.JsonSetterImpl[walletEntities.WalletVal],
	)
	createWallet := walletService.NewWalletCreator(storeCreateWallet)
	getWallets := walletService.NewWalletsGetter(storeImpl.NewWalletsGetter(db.JsonGetterImpl[storeImpl.UserWalletsModel], getWallet))
	getBalance := walletService.NewWalletGetter(getWallet)

	createWalletHandler := wHandlers.NewCreateWalletHandler(deps.TRunner, createWallet)
	getWalletsHandler := wHandlers.NewGetWalletsHandler(deps.TRunner, getWallets)

	// ===== STATIC STORE =====
	staticFileCreator := static_store.NewStaticFileCreatorImpl(deps.Config.StaticDir)
	// staticDirDeleter = static_store.NewStaticDirDeleterImpl(deps.Config.StaticDir)

	// ===== UPLOADER =====
	upld := uploader.NewUploaderFactory(uploader.NewServiceFactory(staticFileCreator))

	// ===== LIMITS =====
	storeGetWithdraws := withdrawsStore.NewWithdrawsGetter(db.JsonGetterImpl[models.Withdraws])
	storeUpdateWithdrawn := withdrawsStore.NewWithdrawUpdater(db.JsonUpdaterImpl[models.WithdrawVal])
	getUpdatedWithdrawn := withdrawsService.NewWithdrawnUpdateGetter(storeGetWithdraws)
	updateWithdrawn := withdrawsService.NewWithdrawnUpdater(getUpdatedWithdrawn, storeUpdateWithdrawn)
	transUpdateWithdrawn := withdrawsService.NewTransWithdrawnUpdater(getWallet, updateWithdrawn)
	getLimits := limits.NewLimitsGetter(storeGetWithdraws, limits.NewLimitsComputer(configurable.LimitedCurrencies))
	getLimit := limits.NewLimitGetter(getWallet, getLimits)
	checkLimit := limits.NewLimitChecker(getLimit)

	// ===== USERS =====
	getUserInfo := userService.NewUserInfoGetter(getWallets, getLimits, deps.Auth.Get)
	getUserInfoHandler := uHandlers.NewGetUserInfoHandler(deps.TRunner, getUserInfo)
	userDetailsCrudEndpoint := users.NewUserDetailsCRUDEndpoint(deps.SimpleDB)

	// ===== HISTORIES =====
	storeGetHistory := histStore.NewHistoryGetter(db.JsonColGetterImpl[histEntities.TransEntry])
	storeStoreTrans := histStore.NewTransStorer(db.JsonSetterImpl[histEntities.TransEntry])
	getHistory := histService.NewHistoryGetter(storeGetHistory)
	storeTrans := histService.NewEntryStorer(getWallet, storeStoreTrans)
	getHistoryHandler := histHandlers.NewGetHistoryHandler(deps.TRunner, getHistory)

	// ===== TRANSACTIONS =====
  walletOwnershipValidator := tValidators.NewWalletOwnershipValidator(getWallet)
	transValidator := tValidators.NewTransactionValidator(checkLimit, tValidators.NewEnoughBalanceValidator(getBalance))
	codeParser := mappers.NewCodeParser(jwtVerifier)
	codeMapper := mappers.NewCodeGenerator(jwtIssuer)
	codeGenerator := service.NewCodeGenerator(walletOwnershipValidator, codeMapper)

	createTrans := tStore.NewTransactionCreator(codeMapper)
	getTrans := tStore.NewTransactionGetter(codeParser)
	transPerformer := tService.NewTransactionPerformer(transUpdateWithdrawn, storeTrans, tService.NewTransBalUpdater(updBal))
	transact := tService.NewTransactionFinalizer(transValidator, transPerformer)
	multiTransact := tService.NewMultiTransactionFinalizer(transact)

	// ===== ATM =====
	atmSecretValidator := atmValidators.NewATMSecretValidator(deps.AtmSecret)
	metaTransByIdValidator := atmValidators.NewMetaTransByIdValidator(getTrans)
	metaTransFromCodeValidator := atmValidators.NewMetaTransFromCodeValidator(codeParser)
	validateWithdrawal := atmValidators.NewWithdrawalValidator(metaTransByIdValidator, transValidator)
	insertedBanknoteValidator := atmValidators.NewInsertedBanknoteValidator(metaTransByIdValidator)
	dispensedBanknoteValidator := atmValidators.NewDispensedBanknoteValidator(metaTransByIdValidator)

	createAtmTrans := atmService.NewATMTransactionCreator(metaTransFromCodeValidator, getUserInfo, createTrans)
	cancelTrans := atmService.NewTransactionCanceler()
	generalFinalizer := atmService.NewGeneralFinalizer(metaTransByIdValidator, multiTransact)
	finalizeDeposit := atmService.NewDepositFinalizer(generalFinalizer)
	finalizeWithdrawal := atmService.NewWithdrawalFinalizer(generalFinalizer)

	atmAuthMiddleware := atmMiddleware.NewATMAuthMiddleware(atmSecretValidator)

	genCodeHandler := atmHandlers.NewGenerateQRCodeHandler(deps.TRunner, codeGenerator)
	createTransHandler := atmHandlers.NewCreateTransactionHandler(deps.TRunner, createAtmTrans)
	onCancelHandler := atmHandlers.NewCancelTransactionHandler(cancelTrans)
	validateWithdrawalHandler := atmHandlers.NewWithdrawalValidationHandler(deps.TRunner, validateWithdrawal)
	completeDepositHandler := atmHandlers.NewCompleteDepostHandler(deps.TRunner, finalizeDeposit)
	completeWithdrawalHandler := atmHandlers.NewCompleteWithdrawalHandler(deps.TRunner, finalizeWithdrawal)
	banknoteEscrowHandler := atmHandlers.NewBanknoteEscrowHandler(deps.TRunner, insertedBanknoteValidator)
	banknoteAcceptedHandler := atmHandlers.NewBanknoteAcceptedHandler(deps.TRunner, insertedBanknoteValidator)
	preBanknoteDispensedHandler := atmHandlers.NewPreBanknoteDispensedHandler(deps.TRunner, dispensedBanknoteValidator)
	postBanknoteDispensedHandler := atmHandlers.NewPostBanknoteDispensedHandler(deps.TRunner, dispensedBanknoteValidator)

	// // ===== TRANSFERS =====
	// convertTransfer := transService.NewTransferConverter(userFromEmail)
	// validateTransfer := transValidators.NewTransferValidator()
	// performTransfer := transService.NewTransferPerformer(multiTransact)
	// transfer := transService.NewTransferer(convertTransfer, validateTransfer, performTransfer)
	// transferHandler := transHandlers.NewTransferHandler(deps.TRunner, transfer)
	//
	// ===== KYC =====
	statusEPFactory := kyc.NewStatusEndpointFactory(deps.SimpleDB)
	passportEndpoint := kyc.NewPassportEndpoint(upld, statusEPFactory)

	return APIDeps{
		Handlers: api.Handlers{
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
				Transfer: func(http.ResponseWriter, *http.Request) {
					panic("unimplemented")
				},
				GetUserInfo: getUserInfoHandler,
				GetHistory:  getHistoryHandler,
				Kyc:         api.KycHandlers{Passport: passportEndpoint},
				UserDetails: userDetailsCrudEndpoint,
				Wallets:     api.WalletHandlers{
					GetAll: getWalletsHandler,
					Create: createWalletHandler,
				},
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
