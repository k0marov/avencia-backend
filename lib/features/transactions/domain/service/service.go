package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

// TODO: having dots in a transactionId (since it is internally a JWT) may result in having dots in url query params, this may lead to bugs


type TransactionIdGetter = func(trans values.MetaTrans) (id string, err error)
type TransactionGetter = func(transactionId string) (values.MetaTrans, error)

type MultiTransactionFinalizer = func(db db.DB, t []values.Transaction) error 
type TransactionFinalizer = func(db db.DB, t values.Transaction) error
type transactionPerformer = func(db db.DB, curBalance core.MoneyAmount, t values.Transaction) error


func NewTransactionIdGetter(genCode mappers.CodeGenerator, genId mappers.TransIdGenerator) TransactionIdGetter {
	return func(trans values.MetaTrans) (id string, err error) {
		code, err := genCode(trans)
		if err != nil {
			return "", core_err.Rethrow("generating code", err)
		}
		return genId(code.Code) , nil
	}
}


func NewTransactionGetter(parseId mappers.TransIdParser, parseCode mappers.CodeParser) TransactionGetter {
	return func(transactionId string) (values.MetaTrans, error) {
		code := parseId(transactionId )
		return parseCode(code)
	}
}

// TODO: not tested 
func NewMultiTransactionFinalizer(finalize TransactionFinalizer) MultiTransactionFinalizer {
	return func(db db.DB, tList []values.Transaction) error {
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
	return func(db db.DB, t values.Transaction) error {
		bal, err := validate(db, t)
		if err != nil {
			return err
		}
		return perform(db, bal, t)
	}
}

func NewTransactionPerformer(updateWithdrawn limitsService.WithdrawnUpdater, addHist histService.TransStorer, updBal walletStore.BalanceUpdater) transactionPerformer {
	return func(db db.DB, curBal core.MoneyAmount, t values.Transaction) error {
		// TODO: maybe move this check inside updateWithdrawn
		if t.Money.Amount.IsNeg() {
			err := updateWithdrawn(db, t)
			if err != nil {
				return core_err.Rethrow("updating withdrawn", err)
			}
		}
		err := addHist(db, t) 
		if err != nil {
			return core_err.Rethrow("adding trans to history", err)
		}

		return updBal(db, t.UserId, t.Money.Currency, curBal.Add(t.Money.Amount))
	}
}
