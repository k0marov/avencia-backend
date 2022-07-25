package service

import (
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
)

type Transferer = func(values.RawTransfer) error
type transferConverter = func(values.RawTransfer) values.Transfer

func NewTransferer(convert transferConverter) Transferer {
	return func(t values.RawTransfer) error {
		panic("unexpected")
	}
}

func NewTransferConverter(userFromEmail auth.UserFromEmail) transferConverter {
	return func(rt values.RawTransfer) values.Transfer {
		panic("unimplemented")
	}
}
