package lib

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/config"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction"
	"github.com/k0marov/avencia-backend/lib/features/limits"
	"github.com/k0marov/avencia-backend/lib/features/transfer"
	"github.com/k0marov/avencia-backend/lib/features/user"
	userService "github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/wallet"
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

// TODO: maybe stop using individual DI integrators for every feature, since it is becoming hard to get individual services from each feature

func Initialize() http.Handler {
	conf := config.LoadConfig()

	fbApp := initFirebase(conf)
	fsClient, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error while initializing firestore client: %v", err)
	}

	authMiddleware := auth.NewAuthMiddleware(fbApp)
	walletServices := wallet.NewWalletServicesImpl(fsClient)
	limitsServices := limits.NewLimitsServicesImpl(fsClient)

	walletDeps := atm_transaction.WalletDeps{
		GetBalance:    walletServices.GetBalance,
		UpdateBalance: walletServices.BalanceUpdater,
	}
	userDeps := atm_transaction.UserDeps{
		GetUserInfo: userService.NewUserInfoGetter(walletServices.GetWallet, limitsServices.GetLimits),
	}
	limitsDeps := atm_transaction.LimitsDeps{
		CheckLimit:          limitsServices.CheckLimit,
		GetUpdatedWithdrawn: limitsServices.GetWithdrawnUpdate,
		UpdateWithdrawn:     limitsServices.UpdateWithdrawn,
	}

	transHandlers := atm_transaction.NewATMTransactionHandlers(conf, fsClient, walletDeps, userDeps, limitsDeps)
	userHandlers := user.NewUserHandlersImpl(userDeps.GetUserInfo)

	transferHandler := transfer.NewTransferHandlerImpl(fsClient, nil, nil) // TODO add userFromEmail and Transact here

	apiRouter := api.NewAPIRouter(api.Handlers{
		GenCode:             transHandlers.GenCode,
		VerifyCode:          transHandlers.VerifyCode,
		CheckBanknote:       transHandlers.CheckBanknote,
		FinalizeTransaction: transHandlers.FinalizeTransaction,
		GetUserInfo:         userHandlers.GetUserInfo,
		Transfer:            transferHandler,
	}, authMiddleware)
	return middleware.Recoverer(apiRouter)
}
