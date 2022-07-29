package service

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade/batch"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

type Transferer = func(values.RawTransfer) error

// transferConverter error may be a ClientError
type transferConverter = func(values.RawTransfer) (values.Transfer, error)

type transferPerformer = func(values.Transfer) error

func NewTransferer(convert transferConverter, validate validators.TransferValidator, perform transferPerformer) Transferer {
	return func(raw values.RawTransfer) error {
		t, err := convert(raw)
		if err != nil {
			return core_err.Rethrow("converting raw transfer data to a transfer", err)
		}
		err = validate(t)
		if err != nil {
			return err
		}
		return perform(t)
	}
}

func NewTransferPerformer(runBatch batch.WriteRunner, transact tService.TransactionFinalizer) transferPerformer {
	return func(t values.Transfer) error {
		return runBatch(func(u fs_facade.BatchUpdater) error {
			// withdraw money from the wallet of caller
			withdrawTrans := transValues.Transaction{
				UserId: t.FromId,
				Money: core.Money{
					Currency: t.Money.Currency,
					Amount:   t.Money.Amount.Neg(),
				},
			}
			err := transact(u, withdrawTrans)
			if err != nil {
				return err
			}
			// deposit money to recipient
			depositTrans := transValues.Transaction{
				UserId: t.ToId,
				Money: core.Money{
					Currency: t.Money.Currency,
					Amount:   t.Money.Amount,
				},
			}
			err = transact(u, depositTrans)
			if err != nil {
				return err
			}
			return nil
		})
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
