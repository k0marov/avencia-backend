package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)


// TODO: change something so that "u fs_facade.BatchUpdater" is not so long 

type TransactionGetter = func(transactionId string) (values.MetaTrans, error)
type TransactionCreator = func(trans values.MetaTrans) (id string, err error)


type TransactionFinalizer = func(u fs_facade.BatchUpdater, t values.Transaction) error
type transactionPerformer = func(u fs_facade.BatchUpdater, curBalance core.MoneyAmount, t values.Transaction) error

func NewTransactionFinalizer(validate validators.TransactionValidator, perform transactionPerformer) TransactionFinalizer {
	return func(u fs_facade.BatchUpdater, t values.Transaction) error {
		bal, err := validate(t)
		if err != nil {
			return err
		}
		return perform(u, bal, t)
	}
}

// TODO: rename fs_facade to db_facade 

func NewTransactionPerformer(updateWithdrawn limitsService.WithdrawnUpdater, addHist histService.TransStorer, updBal walletStore.BalanceUpdater) transactionPerformer {
	return func(u fs_facade.BatchUpdater, curBal core.MoneyAmount, t values.Transaction) error {
		if t.Money.Amount.IsNeg() {
			err := updateWithdrawn(u, t)
			if err != nil {
				return core_err.Rethrow("updating withdrawn", err)
			}
		}
		err := addHist(u, t) 
		if err != nil {
			return core_err.Rethrow("adding trans to history", err)
		}

		return updBal(u, t.UserId, t.Money.Currency, curBal.Add(t.Money.Amount))
	}
}
