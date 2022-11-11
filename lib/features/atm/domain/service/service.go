package service

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	tService "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	tStore "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/store"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	uService "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/service"
)

type ATMTransactionCreator = func(db.TDB, values.TransFromQRCode) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error

type DepositFinalizer = func(db.TDB, values.DepositData) error
type WithdrawalFinalizer = func(db.TDB, values.WithdrawalData) error

// TODO: add validation that there is no active transaction for this user
func NewATMTransactionCreator(
	val validators.MetaTransByCodeValidator, 
	getUser uService.UserInfoGetter,
	create tStore.TransactionCreator, 
) ATMTransactionCreator {
	return func(db db.TDB, nt values.TransFromQRCode) (values.CreatedTransaction, error) {
		metaTrans, err := val(nt.QRCodeText, nt.Type)
		if err != nil {
			return values.CreatedTransaction{}, err
		}
		user, err := getUser(db, metaTrans.CallerId)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting user info", err)
		}
		transId, err := create(metaTrans)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting the transaction id", err)
		}
		return values.CreatedTransaction{Id: transId, UserInfo: user}, nil
	}
}

// TODO: here the user's current transaction may be reset
// TODO: invalidate the transactionId
func NewTransactionCanceler() TransactionCanceler {
	return func(id string) error {
		return nil
	}
}

type generalFinalizer = func(db db.TDB, transId string, wantType tValues.TransactionType, m []core.Money) error

func NewDepositFinalizer(generalFinalizer generalFinalizer) DepositFinalizer {
	return func(db db.TDB, dd values.DepositData) error {
		return generalFinalizer(db, dd.TransactionId, tValues.Deposit, dd.Received)
	}
}

func NewWithdrawalFinalizer(generalFinalizer generalFinalizer) WithdrawalFinalizer {
	return func(db db.TDB, wd values.WithdrawalData) error {
		return generalFinalizer(db, wd.TransactionId, tValues.Withdrawal, []core.Money{wd.Money})
	}
}

// TODO: here the user's current transaction may be reset
// TODO: invalidate the transactionId
func NewGeneralFinalizer(validate validators.MetaTransByIdValidator, finalize tService.MultiTransactionFinalizer) generalFinalizer {
	return func(db db.TDB, transId string, tType tValues.TransactionType, m []core.Money) error {
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
				WalletId: metaTrans.WalletId,
				Money:  m.Amount,
			})
		}
		return finalize(db, t)
	}
}

