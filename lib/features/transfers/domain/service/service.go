package service

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	authStore "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

type Transferer = func(transactionalDB db.TDB, t values.RawTransfer) error

type transferConverter = func(db.TDB, values.RawTransfer) (values.Transfer, error)
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
// 			UserId: t.FromId,
// 			Money: core.Money{
// 				Currency: t.Money.Currency,
// 				Amount:   t.Money.Amount.Neg(),
// 			},
// 		}
// 		depositTrans := transValues.Transaction{
// 			Source: transValues.TransSource{
// 				Type:   transValues.Transfer,
// 				Detail: t.FromId,
// 			},
// 			UserId: t.ToId,
// 			Money: core.Money{
// 				Currency: t.Money.Currency,
// 				Amount:   t.Money.Amount,
// 			},
// 		}
// 		return transact(db, []transValues.Transaction{withdrawTrans, depositTrans})
// 	}
// }

func NewTransferConverter(userFromEmail authStore.UserByEmailGetter, getWallet wService.WalletGetter) transferConverter {
	return func(db db.TDB, t values.RawTransfer) (values.Transfer, error) {
		user, err := userFromEmail(t.ToEmail)
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
		return values.Transfer{
			FromId: t.FromId,
			SourceWallet: wallet,
			ToId:   user.Id,
			Amount:  t.Amount,
		}, nil
	}
}
