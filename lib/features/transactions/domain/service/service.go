package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

type TransactionFinalizer = func(u firestore_facade.BatchUpdater, t values.Transaction) error
type transactionPerformer = func(u firestore_facade.BatchUpdater, curBalance core.MoneyAmount, t values.Transaction) error

func NewTransactionFinalizer(validate validators.TransactionValidator, perform transactionPerformer) TransactionFinalizer {
	return func(u firestore_facade.BatchUpdater, t values.Transaction) error {
		bal, err := validate(t)
		if err != nil {
			return err
		}
		return perform(u, bal, t)
	}
}

func NewTransactionPerformer(updBal walletStore.BalanceUpdater, updateWithdrawn limitsService.WithdrawUpdater) transactionPerformer {
	return func(u firestore_facade.BatchUpdater, curBal core.MoneyAmount, t values.Transaction) error {
		if t.Money.Amount.IsNeg() {
			err := updateWithdrawn(u, t)
			if err != nil {
				return core_err.Rethrow("updating withdrawn", err)
			}
		}
		return updBal(u, t.UserId, t.Money.Currency, curBal.Add(t.Money.Amount))
	}
}
