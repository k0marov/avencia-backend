package service

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

type Transferer = func(values.RawTransfer) error

// transferConverter error may be a ClientError
type transferConverter = func(values.RawTransfer) (values.Transfer, error)

// transferValidator error may be a ClientError
type transferValidator = func(values.Transfer) error

// TODO: try to simplify
func NewTransferer(convert transferConverter, validate transferValidator, runBatch batch.WriteRunner, transact atmService.TransactionFinalizer) Transferer {
	return func(raw values.RawTransfer) error {
		t, err := convert(raw)
		if err != nil {
			return core_err.Rethrow("converting raw transfers data to a transfers", err)
		}
		err = validate(t)
		if err != nil {
			return err
		}
		return runBatch(func(u firestore_facade.BatchUpdater) error {
			// withdraw money from the wallets of caller
			withdrawTrans := transValues.Transaction{
				UserId: t.FromId,
				Money: core.Money{
					Currency: t.Money.Currency,
					Amount:   t.Money.Amount.Neg(),
				},
			}
			err = transact(u, withdrawTrans)
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

func NewTransferValidator() transferValidator {
	return func(t values.Transfer) error {
		if t.Money.Amount.IsNeg() {
			return client_errors.NegativeTransferAmount
		}
		if t.Money.Amount.IsEqual(core.NewMoneyAmount(0)) {
			return client_errors.TransferingZero
		}
		if t.ToId == t.FromId {
			return client_errors.TransferingToYourself

		}
		return nil
	}
}
