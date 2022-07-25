package service

import (
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	"github.com/k0marov/avencia-backend/lib/core/core_errors"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
)

type Transferer = func(values.RawTransfer) error

// transferConverter error may be a ClientError
type transferConverter = func(values.RawTransfer) (values.Transfer, error)

// TODO: try to simplify
func NewTransferer(convert transferConverter, runBatch batch.WriteRunner, transact atmService.TransactionFinalizer) Transferer {
	return func(raw values.RawTransfer) error {
		if raw.Money.Amount < 0 {
			return client_errors.NegativeTransferAmount
		}
		t, err := convert(raw)
		if err != nil {
			return fmt.Errorf("converting raw transfer data to a transfer: %w", err)
		}

		return runBatch(func(u firestore_facade.BatchUpdater) error {
			// withdraw money from the wallet of caller
			withdrawTrans := transValues.Transaction{
				UserId: t.FromId,
				Money: core.Money{
					Currency: t.Money.Currency,
					Amount:   -t.Money.Amount,
				},
			}
			err = transact(u, withdrawTrans)
			if err != nil {
				return fmt.Errorf("withdrawing money from caller: %w", err)
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
				return fmt.Errorf("depositing money to recipeint: %w", err)
			}
			return nil
		})
	}
}

func NewTransferConverter(userFromEmail auth.UserFromEmail) transferConverter {
	return func(t values.RawTransfer) (values.Transfer, error) {
		user, err := userFromEmail(t.ToEmail)
		if err == core_errors.ErrNotFound {
			return values.Transfer{}, client_errors.NotFound
		}
		if err != nil {
			return values.Transfer{}, fmt.Errorf("while getting transfer recepient from its email: %w", err)
		}
		return values.Transfer{
			FromId: t.FromId,
			ToId:   user.Id,
			Money:  t.Money,
		}, nil
	}
}
