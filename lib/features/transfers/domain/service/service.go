package service

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type Transferer = func(transactionalDB db.TDB, t values.RawTransfer) error

type transferConverter = func(db.TDB, values.RawTransfer) (values.Transfer, error)
type walletFinder = func(db db.TDB, userId string, needCurrency core.Currency) (wEntities.Wallet, error)
type transferPerformer = func(transactionalDB db.TDB, t values.Transfer) error

// func NewTransferer(convert transferConverter, validate validators.TransferValidator, perform transferPerformer) Transferer {
// 	return func(db db.TDB, raw values.RawTransfer) error {
// 		t, err := convert(raw)
// 		if err != nil {
// 			return core_err.Rethrow("converting raw transfer data to a transfer", err)
// 		}
// 		err = validate(t)
// 		if err != nil {
// 			return err
// 		}
// 		return perform(db, t)
// 	}
// }
//
// func NewTransferPerformer(transact tService.MultiTransactionFinalizer) transferPerformer {
// 	return func(db db.TDB, t values.Transfer) error {
// 		withdrawTrans := transValues.Transaction{
// 			Source: transValues.TransSource{
// 				Type:   transValues.Transfer,
// 				Detail: t.ToId,
// 			},
// 			WalletId: ,
// 			Money:   t.Amount.Neg(),
// 		}
// 		depositTrans := transValues.Transaction{
// 			Source: transValues.TransSource{
// 				Type:   transValues.Transfer,
// 				Detail: t.FromId,
// 			},
// 			UserId: t.ToId,
// 			Money:   t.Amount,
// 		}
// 		return transact(db, []transValues.Transaction{withdrawTrans, depositTrans})
// 	}
// }

func NewWalletFinder(getWallets wService.WalletsGetter) walletFinder {
	return func(db db.TDB, userId string, needCurrency core.Currency) (wEntities.Wallet, error) {
    wallets, err := getWallets(db, userId) 
    if err != nil {
    	return wEntities.Wallet{}, core_err.Rethrow("getting recipient's wallets", err)
    }
    for _, w := range wallets {
    	if w.Currency == needCurrency {
    		return w, nil
    	}
    }
    return wEntities.Wallet{}, client_errors.ProperWalletNotFound
	}
}

func NewTransferConverter(userFromEmail authStore.UserByEmailGetter, getWallet wService.WalletGetter, findWallet walletFinder) transferConverter {
	return func(db db.TDB, t values.RawTransfer) (values.Transfer, error) {
		toUser, err := userFromEmail(t.ToEmail)
		if core_err.IsNotFound(err) {
			return values.Transfer{}, client_errors.NotFound
		}
		if err != nil {
			return values.Transfer{}, core_err.Rethrow("while getting transfers recepient from its email", err)
		}
		wallet, err := getWallet(db, t.SourceWalletId)
		if err != nil {
			return values.Transfer{}, core_err.Rethrow("getting source wallet", err)
		}
		toWallet, err := findWallet(db, toUser.Id, wallet.Currency)
		if err != nil {
			return values.Transfer{}, core_err.Rethrow("finding a fitting target wallet", err)
		}
		return values.Transfer{
			FromId: t.FromId,
			FromWallet: wallet,
			ToWallet: toWallet,
			ToId:   toUser.Id,
			Amount:  t.Amount,
		}, nil
	}
}
