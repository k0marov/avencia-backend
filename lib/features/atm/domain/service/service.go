package service

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type ATMTransactionCreator = func(values.TransFromQRCode) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error

type DepositFinalizer = func(db.DB, values.DepositData) error
type WithdrawalFinalizer = func(db.DB, values.WithdrawalData) error

// TODO: add validation that there is no active transaction for this user
func NewATMTransactionCreator(validate validators.MetaTransByCodeValidator, createTrans tService.TransactionIdGetter) ATMTransactionCreator {
	return func(nt values.TransFromQRCode) (values.CreatedTransaction, error) {
		metaTrans, err := validate(nt.QRCodeText, nt.Type)
		if err != nil {
			return values.CreatedTransaction{}, err
		}
		transId, err := createTrans(metaTrans)
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

type generalFinalizer = func(db db.DB, transId string, wantType tValues.TransactionType, m []core.Money) error

func NewDepositFinalizer(generalFinalizer generalFinalizer) DepositFinalizer {
	return func(db db.DB, dd values.DepositData) error {
		return generalFinalizer(db, dd.TransactionId, tValues.Deposit, dd.Received)
	}
}

func NewWithdrawalFinalizer(generalFinalizer generalFinalizer) WithdrawalFinalizer {
	return func(db db.DB, wd values.WithdrawalData) error {
		return generalFinalizer(db, wd.TransactionId, tValues.Withdrawal, []core.Money{wd.Money})
	}
}

func NewGeneralFinalizer(validate validators.MetaTransByIdValidator, finalize tService.MultiTransactionFinalizer) generalFinalizer {
	return func(db db.DB, transId string, tType tValues.TransactionType, m []core.Money) error {
		metaTrans, err := validate(transId, tType)
		if err != nil {
			return err
		}

		source := tValues.TransSource{
			Type:   tValues.Cash,
			Detail: "",
		}

		var t []tValues.Transaction
		for _, m := range m {
			t = append(t, tValues.Transaction{
				Source: source,
				UserId: metaTrans.UserId,
				Money:  m,
			})
		}
		return finalize(db, t)
	}
}

