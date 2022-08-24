package service

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/db/firestore"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

type DeliveryTransferer = func(t values.RawTransfer) error 

func NewDeliveryTransferer(runT firestore.TransactionRunner, transfer Transferer) DeliveryTransferer {
	return func(t values.RawTransfer) error {
		return runT(func(db db.DB) error {
			return transfer(db, t)
		})
	}
}

type Transferer = func(transactionalDB db.DB, t values.RawTransfer) error

type transferConverter = func(values.RawTransfer) (values.Transfer, error)
type transferPerformer = func(transactionalDB db.DB, t values.Transfer) error

func NewTransferer(convert transferConverter, validate validators.TransferValidator, perform transferPerformer) Transferer {
	return func(db db.DB, raw values.RawTransfer) error {
		t, err := convert(raw)
		if err != nil {
			return core_err.Rethrow("converting raw transfer data to a transfer", err)
		}
		err = validate(t)
		if err != nil {
			return err
		}
		return perform(db, t)
	}
}

func NewTransferPerformer(transact tService.MultiTransactionFinalizer) transferPerformer {
	return func(db db.DB, t values.Transfer) error {
		withdrawTrans := transValues.Transaction{
			Source: transValues.TransSource{
				Type:   transValues.Transfer,
				Detail: t.ToId,
			},
			UserId: t.FromId,
			Money: core.Money{
				Currency: t.Money.Currency,
				Amount:   t.Money.Amount.Neg(),
			},
		}
		depositTrans := transValues.Transaction{
			Source: transValues.TransSource{
				Type:   transValues.Transfer,
				Detail: t.FromId,
			},
			UserId: t.ToId,
			Money: core.Money{
				Currency: t.Money.Currency,
				Amount:   t.Money.Amount,
			},
		}
		return transact(db, []transValues.Transaction{withdrawTrans, depositTrans})
	}
}

func NewTransferConverter(userFromEmail auth.UserFromEmail) transferConverter {
	return func(t values.RawTransfer) (values.Transfer, error) {
		user, err := userFromEmail(t.ToEmail)
		if err == core_err.ErrNotFound {
			return values.Transfer{}, client_errors.NotFound
		}
		if err != nil {
			return values.Transfer{}, core_err.Rethrow("while getting transfers recepient from its email", err)
		}
		return values.Transfer{
			FromId: t.FromId,
			ToId:   user.Id,
			Money:  t.Money,
		}, nil
	}
}
