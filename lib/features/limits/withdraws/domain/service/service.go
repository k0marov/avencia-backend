package service

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/store"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	wService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
)

type WithdrawnUpdater = func(db db.TDB, userId string, tMoney core.Money) error
type TransWithdrawnUpdater = func(db db.TDB, t tValues.Transaction) error

type withdrawnUpdateGetter = func(db db.TDB, userId string, tMoney core.Money) (core.Money, error)

func NewWithdrawnUpdater(getValue withdrawnUpdateGetter, update store.WithdrawUpdater) WithdrawnUpdater {
	return func(db db.TDB, userId string, tMoney core.Money) error {
		if tMoney.Amount.IsPos() { // it is a deposit - no update needed
			return nil
		}
		newWithdrawn, err := getValue(db, userId, tMoney)
		if err != nil {
			return core_err.Rethrow("getting new withdrawn value", err)
		}
		return update(db, userId, newWithdrawn)
	}
}

func NewTransWithdrawnUpdater(getWallet wService.WalletGetter, upd WithdrawnUpdater) TransWithdrawnUpdater {
	return func(db db.TDB, t tValues.Transaction) error {
		wallet, err := getWallet(db, t.WalletId)
		if err != nil {
			return core_err.Rethrow("getting the wallet", err)
		}
		return upd(db, wallet.OwnerId, t.Money)
	}
}

func NewWithdrawnUpdateGetter(getWithdraws store.WithdrawsGetter) withdrawnUpdateGetter {
	return func(db db.TDB, userId string, tMoney core.Money) (core.Money, error) {
		withdraws, err := getWithdraws(db, userId)
		if err != nil {
			return core.Money{}, core_err.Rethrow("getting limits", err)
		}
		curWithdrawn := withdraws[tMoney.Currency]
		newWithdraw := tMoney.Amount.Neg()

		result := core.Money{Currency: tMoney.Currency}
		if configurable.IsWithdrawLimitRelevant(time.Unix(curWithdrawn.UpdatedAt, 0)) {
			result.Amount = curWithdrawn.Withdrawn.Add(newWithdraw)
		} else {
			result.Amount = newWithdraw
		}
		return result, nil
	}
}
