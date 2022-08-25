package service

import (
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

type DeliveryTransferer = func(t values.RawTransfer) error 

func NewDeliveryTransferer(runT db.TransactionRunner, transfer Transferer) DeliveryTransferer {
	return func(t values.RawTransfer) error {
		return runT(func(db db.DB) error {
			return transfer(db, t)
		})
	}
}

