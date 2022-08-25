package service

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tMappers "github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type ATMTransactionCreator = func(values.NewTrans) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error

type DepositFinalizer = func(db.DB, values.DepositData) error
type WithdrawalFinalizer = func(db.DB, values.WithdrawalData) error

// TODO: add validation that there is no active transaction for this user
func NewATMTransactionCreator(parseCode tMappers.CodeParser, createTrans tService.TransactionIdGetter) ATMTransactionCreator {
	return func(nt values.NewTrans) (values.CreatedTransaction, error) {
		trans, err := parseCode(nt.QRCodeText)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting transaction from qr code", err)
		}
		// TODO: maybe move this to a separate validator
		if trans.Type != nt.Type {
			return values.CreatedTransaction{}, client_errors.InvalidTransactionType
		}
		transId, err := createTrans(trans)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting the transaction id", err)
		}
		return values.CreatedTransaction{Id: transId}, nil
	}
}

// TODO: here the user's current transaction may be reset
// TODO: invalidate the transactionId
func NewTransactionCanceler() TransactionCanceler {
	return func(id string) error {
		return nil
	}
}

// TODO: somehow DRY getting metaTrans and validating metaTrans.Type from ATMTransactionCreator, DepositFinalizer, WithdrawalFinalizer. The new function (validator) could also be used in banknote validation services


func NewDepositFinalizer(getTrans tService.TransactionGetter, finalize tService.MultiTransactionFinalizer) DepositFinalizer {
	return func(db db.DB, dd values.DepositData) error {
		metaTrans, err := getTrans(dd.TransactionId)
		if err != nil {
			return core_err.Rethrow("getting transaction from trans id", err)
		}
		if metaTrans.Type != tValues.Deposit {
			return client_errors.InvalidTransactionType
		}

		source := tValues.TransSource{
			Type:   tValues.Cash,
			Detail: "",
		}

		var t []tValues.Transaction
		for _, m := range dd.Received {
			t = append(t, tValues.Transaction{
				Source: source,
				UserId: metaTrans.UserId,
				Money:  m,
			})
		}
		return finalize(db, t)
	}
}

func NewWithdrawalFinalizer(getTrans tService.TransactionGetter, finalize tService.TransactionFinalizer) WithdrawalFinalizer {
	return func(db db.DB, wd values.WithdrawalData) error {
		metaTrans, err := getTrans(wd.TransactionId)
		if err != nil {
			return core_err.Rethrow("getting transaction from trans id", err)
		}
		if metaTrans.Type != tValues.Withdrawal {
			return client_errors.InvalidTransactionType
		}

		t := tValues.Transaction{
			Source: tValues.TransSource{
				Type: tValues.Cash, 
				Detail: "", 
			},
			UserId: metaTrans.UserId,
			Money:  wd.Money,
		}
		return finalize(db, t)
	}
}
