package service

import (
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_errors"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
)

type Transferer = func(values.RawTransfer) error

// transferConverter error may be a ClientError
type transferConverter = func(values.RawTransfer) (values.Transfer, error)

func NewTransferer(convert transferConverter, transact atmService.TransactionFinalizer) Transferer {
	return func(raw values.RawTransfer) error {
		if raw.Money.Amount < 0 {
			return client_errors.NegativeTransferAmount
		}
		_, err := convert(raw)
		if err != nil {
			return fmt.Errorf("while converting raw transfer data to a transfer: %w", err)
		}
		return nil
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
