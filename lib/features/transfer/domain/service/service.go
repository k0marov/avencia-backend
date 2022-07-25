package service

import "github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"

type Transferer = func(values.RawTransfer) error

func NewTransferer() Transferer {
	return func(t values.RawTransfer) error {
		panic("unexpected")
	}
}
