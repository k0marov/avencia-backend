package service

import (
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
)

type DeliveryHistoryGetter = func(userId string) ([]entities.TransEntry, error)

func NewDeliveryHistoryGetter(simpleDB db.DB, getHistory HistoryGetter) DeliveryHistoryGetter {
	return func(userId string) ([]entities.TransEntry, error) {
		return getHistory(simpleDB, userId)
	}
}
