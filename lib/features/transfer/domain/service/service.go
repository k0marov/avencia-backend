package service

import (
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_errors"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
)

type Transferer = func(values.RawTransfer) error

// transferConverter error may be a ClientError
type transferConverter = func(values.RawTransfer) (values.Transfer, error)

func NewTransferer(convert transferConverter) Transferer {
	return func(t values.RawTransfer) error {
		panic("unimplemented")
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
