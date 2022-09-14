package service_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestTransactionIdGetter(t *testing.T) {
	trans := RandomMetaTrans() 
	code := RandomGeneratedCode()
	id := RandomString() 
	t.Run("error case - generating code throws", func(t *testing.T) {
		genCode := func(gotTrans values.MetaTrans) (values.GeneratedCode, error) {
			if gotTrans == trans {
				return values.GeneratedCode{}, RandomError() 
			}
			panic("unexpected")
		}
		_, err := service.NewTransactionIdGetter(genCode, nil)(trans) 
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		genCode := func(values.MetaTrans) (values.GeneratedCode, error) {
			return 	code, nil
		}
		genId := func(gotCode string) string {
			if gotCode == code.Code {
				return id
			}
			panic("unexpected")
		}

		gotId, err := service.NewTransactionIdGetter(genCode, genId)(trans) 
		AssertNoError(t, err)
		Assert(t, gotId, id, "returned id")
	})
}

func TestTransactionGetter(t *testing.T) {
	id := RandomString() 
	code := RandomString()
	trans := RandomMetaTrans() 

	parseId := func(gotId string) string {
		if gotId == id {
			return code 
		}
		panic("unexpected")
	}

	t.Run("should forward to parseCode", func(t *testing.T) {
		tErr := RandomError() 
		parseCode := func(gotCode string) (values.MetaTrans, error) {
			return trans, tErr
		}

		gotTrans, err := service.NewTransactionGetter(parseId, parseCode)(id)
		AssertError(t, err, tErr)
		Assert(t, gotTrans, trans, "returned transaction")
	})


}



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
	wantNewBal := curBalance.Add(trans.Money.Amount)

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
		if gotDB == mockDB && gotTrans == trans {
			return nil
		}
		panic("unexpected")
	}
	t.Run("adding transaction to history throws", func(t *testing.T) {
		addHist := func(db.DB, values.Transaction) error {
			return RandomError()
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, nil)(mockDB, curBalance, trans)
		AssertSomeError(t, err)
	})

	updBal := func(gotDB db.DB, user string, newBal core.Money) error {
		if gotDB == mockDB && 
		   user == trans.UserId && 
		   newBal.Currency == trans.Money.Currency && 
		   newBal.Amount.IsEqual(wantNewBal) {
			return nil
		}
		panic("unexpected")
	}
	t.Run("updating balance throws", func(t *testing.T) {
		updBal := func(db.DB, string, core.Money) error {
			return RandomError()
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(mockDB, curBalance, trans)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(mockDB, curBalance, trans)
		AssertNoError(t, err)
	})
}
