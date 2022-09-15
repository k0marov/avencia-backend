package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestTransactionFinalizer(t *testing.T) {
	transaction := RandomTransactionData()
	mockDB := NewStubDB()
	t.Run("error case - validation throws", func(t *testing.T) {
		err := RandomError()
		validate := func(gotDB db.DB, t values.Transaction) (core.MoneyAmount, error) {
			if gotDB == mockDB && t == transaction {
				return core.NewMoneyAmount(0), err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransactionFinalizer(validate, nil)(mockDB, transaction)
		AssertError(t, gotErr, err)
	})
	t.Run("forward case - return whatever performTransaction returns", func(t *testing.T) {
		wantErr := RandomError()
		currentBalance := RandomPosMoneyAmount()
		validate := func(db.DB, values.Transaction) (core.MoneyAmount, error) {
			return currentBalance, nil
		}
		performTransaction := func(gotDB db.DB, curBal core.MoneyAmount, trans values.Transaction) error {
			if curBal == currentBalance && trans == transaction {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(validate, performTransaction)(mockDB, transaction)
		AssertError(t, err, wantErr)
	})
}

func TestTransactionPerformer(t *testing.T) {
	mockDB := NewStubDB()
	curBalance := RandomPosMoneyAmount()
	trans := RandomTransactionData()
	// wantNewBal := core.Money{
	// 	Currency: trans.Money.Currency,
	// 	Amount:   curBalance.Add(trans.Money.Amount),
	// }

	updateWithdrawn := func(db.DB, values.Transaction) error {
		return nil
	}
	t.Run("updating withdrawn throws", func(t *testing.T) {
		updateWithdrawn := func(gotDB db.DB, gotTrans values.Transaction) error {
			if gotDB == mockDB && gotTrans == trans {
				return RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(updateWithdrawn, nil, nil)(mockDB, curBalance, trans)
		AssertSomeError(t, err)
	})

	addHist := func(gotDB db.DB, gotTrans values.Transaction) error {
		return nil
	}
	t.Run("adding transaction to history throws", func(t *testing.T) {
		addHist := func(gotDB db.DB, gotTrans values.Transaction) error {
			if gotDB == mockDB && gotTrans == trans {
				return RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, nil)(mockDB, curBalance, trans)
		AssertSomeError(t, err)
	})

	t.Run("forward case - update balance", func(t *testing.T) {
		tErr := RandomError()
		updBal := func(gotDB db.DB, curBal core.MoneyAmount, gotTrans values.Transaction) error {
			if gotDB == mockDB && curBal == curBalance && gotTrans == trans {
				return tErr
			}
			panic("unexpected")
		}
		// updBal := func(gotDB db.DB, user string, newBal core.Money) error {
		// 	if gotDB == mockDB && user == trans.UserId && newBal.Currency == trans.Money.Currency && newBal.Amount.IsEqual(wantNewBal.Amount) { return tErr
		// 	}
		// 	panic("unexpected")
		// }
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(mockDB, curBalance, trans)
		AssertError(t, err, tErr)
	})
}

func TestTransBalUpdater(t *testing.T) {

	mockDB := NewStubDB()
	curBalance := RandomPosMoneyAmount()
	trans := RandomTransactionData()
	wantNewBal := curBalance.Add(trans.Money.Amount)

	tErr := RandomError()
	updBal := func(gotDB db.DB, user string, newBal core.Money) error {
		if gotDB == mockDB && user == trans.UserId && newBal.Currency == trans.Money.Currency && newBal.Amount.IsEqual(wantNewBal) {
			return tErr
		}
		panic("unexpected")
	}
	err := service.NewTransBalUpdater(updBal)(mockDB, curBalance, trans) 
	AssertError(t, err, tErr)
}
