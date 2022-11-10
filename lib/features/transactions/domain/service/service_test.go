package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestCodeGenerator(t *testing.T) {
	mockDB := NewStubDB()
	metaTrans := RandomMetaTrans()
	ownerValidator := func(gotDB db.TDB, gotTrans values.MetaTrans) error {
		if gotDB == mockDB && gotTrans == metaTrans {
			return nil
		}
		panic("unexpected")
	}
	t.Run("error case - owner validation throws", func(t *testing.T) {
		tErr := RandomError()
     ownerValidator := func(db.TDB, values.MetaTrans) error {
     	 return tErr
     }
     _, err := service.NewCodeGenerator(ownerValidator, nil)(mockDB, metaTrans) 
     AssertError(t, err, tErr)
	})
	t.Run("forward case - forward to code generator mapper", func(t *testing.T) {
		tCode := RandomGeneratedCode()
		tErr := RandomError()
		mapper := func(gotTrans values.MetaTrans) (values.GeneratedCode, error) {
			if gotTrans == metaTrans {
				return tCode, tErr
			}
			panic("unexpected")
		}
		gotCode, gotErr := service.NewCodeGenerator(ownerValidator, mapper)(mockDB, metaTrans) 
		AssertError(t, gotErr, tErr)
		Assert(t, gotCode, tCode, "returned code")
	})
}

func TestMultiTransactionFinalizer(t *testing.T) {
	ts := []values.Transaction{
		RandomTransactionData(),
		RandomTransactionData(),
		RandomTransactionData(),
	}
	t.Run("error case - one of the transactions fails", func(t *testing.T) {
		finalize := func(db.TDB, values.Transaction) error {
			return RandomError()
		}
		err := service.NewMultiTransactionFinalizer(finalize)(NewStubDB(), ts)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		called := []values.Transaction{}
		finalize := func(gotDB db.TDB, gotT values.Transaction) error {
			called = append(called, gotT)
			return nil
		}
		err := service.NewMultiTransactionFinalizer(finalize)(NewStubDB(), ts)
		AssertNoError(t, err)
		Assert(t, called, ts, "array of finalized transactions")
	})
}

func TestTransactionFinalizer(t *testing.T) {
	transaction := RandomTransactionData()
	mockDB := NewStubDB()
	t.Run("error case - validation throws", func(t *testing.T) {
		err := RandomError()
		validate := func(gotDB db.TDB, t values.Transaction) (core.MoneyAmount, error) {
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
		validate := func(db.TDB, values.Transaction) (core.MoneyAmount, error) {
			return currentBalance, nil
		}
		performTransaction := func(gotDB db.TDB, curBal core.MoneyAmount, trans values.Transaction) error {
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

	updateWithdrawn := func(db.TDB, values.Transaction) error {
		return nil
	}
	t.Run("updating withdrawn throws", func(t *testing.T) {
		updateWithdrawn := func(gotDB db.TDB, gotTrans values.Transaction) error {
			if gotDB == mockDB && gotTrans == trans {
				return RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(updateWithdrawn, nil, nil)(mockDB, curBalance, trans)
		AssertSomeError(t, err)
	})

	addHist := func(gotDB db.TDB, gotTrans values.Transaction) error {
		return nil
	}
	t.Run("adding transaction to history throws", func(t *testing.T) {
		addHist := func(gotDB db.TDB, gotTrans values.Transaction) error {
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
		updBal := func(gotDB db.TDB, curBal core.MoneyAmount, gotTrans values.Transaction) error {
			if gotDB == mockDB && curBal == curBalance && gotTrans == trans {
				return tErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(mockDB, curBalance, trans)
		AssertError(t, err, tErr)
	})
}

func TestTransBalUpdater(t *testing.T) {
	mockDB := NewStubDB()
	curBalance := RandomPosMoneyAmount()
	trans := RandomTransactionData()
	wantNewBal := curBalance.Add(trans.Money)

	tErr := RandomError()
	updBal := func(gotDB db.TDB, walletId string, newBal core.MoneyAmount) error {
		if gotDB == mockDB && newBal.IsEqual(wantNewBal) {
			return tErr
		}
		panic("unexpected")
	}
	err := service.NewTransBalUpdater(updBal)(mockDB, curBalance, trans)
	AssertError(t, err, tErr)
}
