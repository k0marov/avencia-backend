package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	histService "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/service"
	withdrawsService "github.com/AvenciaLab/avencia-backend/lib/features/limits/withdraws/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)




type MultiTransactionFinalizer = func(db db.TDB, t []values.Transaction) error 
type TransactionFinalizer = func(db db.TDB, t values.Transaction) error
type transactionPerformer = func(db db.TDB, curBalance core.MoneyAmount, t values.Transaction) error


func NewMultiTransactionFinalizer(finalize TransactionFinalizer) MultiTransactionFinalizer {
	return func(db db.TDB, tList []values.Transaction) error {
		for _, t := range tList {
			err := finalize(db, t)
			if err != nil {
				return core_err.Rethrow("finalizing one of the transactions", err)
			}
		}
		return nil 
	}
}


func NewTransactionFinalizer(validate validators.TransactionValidator, perform transactionPerformer) TransactionFinalizer {
	return func(db db.TDB, t values.Transaction) error {
		bal, err := validate(db, t)
		if err != nil {
			return err
		}
		return perform(db, bal, t)
	}
}

type transBalUpdater = func(db db.TDB, curBal core.MoneyAmount, t values.Transaction) error

func NewTransactionPerformer(updWithdrawn withdrawsService.WithdrawnUpdater, addHist histService.TransStorer, updBal transBalUpdater) transactionPerformer {
	return func(db db.TDB, curBal core.MoneyAmount, t values.Transaction) error {
		err := updWithdrawn(db, t)
		if err != nil {
			return core_err.Rethrow("updating withdrawn", err)
		}
		err = addHist(db, t) 
		if err != nil {
			return core_err.Rethrow("adding trans to history", err)
		}

		return updBal(db, curBal, t)
	}
}

func NewTransBalUpdater(updBal store.BalanceUpdater) transBalUpdater {
	return func(db db.TDB, curBal core.MoneyAmount, t values.Transaction) error {
		return updBal(db, t.UserId, core.Money{
			Currency: t.Money.Currency,
			Amount:   curBal.Add(t.Money.Amount),
		})
	}
}


