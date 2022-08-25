package service_test

import (
	"reflect"
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func TestATMTransactionCreator(t *testing.T) {
	newTrans := values.NewTrans{
		Type:       RandomTransactionType(),
		QRCodeText: RandomString(),
	}
	metaTrans := tValues.MetaTrans{
		Type:   newTrans.Type,
		UserId: RandomString(),
	}
	id := RandomString()
	getTrans := func(code string) (tValues.MetaTrans, error) {
		return metaTrans, nil
	}
	t.Run("error case - parsing qr code throws", func(t *testing.T) {
		getTrans := func(code string) (tValues.MetaTrans, error) {
			if code == newTrans.QRCodeText {
				return tValues.MetaTrans{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewATMTransactionCreator(getTrans, nil)(newTrans)
		AssertSomeError(t, err)
	})
	t.Run("error case - transaction type is not right", func(t *testing.T) {
		var wrongMetaTrans tValues.MetaTrans
		if newTrans.Type == tValues.Deposit {
			wrongMetaTrans.Type = tValues.Withdrawal
		} else {
			wrongMetaTrans.Type = tValues.Deposit
		}
		getTrans := func(string) (tValues.MetaTrans, error) {
			return wrongMetaTrans, nil
		}
		_, err := service.NewATMTransactionCreator(getTrans, nil)(newTrans)
		AssertError(t, err, client_errors.InvalidTransactionType)
	})

	getId := func(tValues.MetaTrans) (string, error) {
		return id, nil
	}

	t.Run("error case - getting transaction id throws", func(t *testing.T) {
		getId := func(trans tValues.MetaTrans) (string, error) {
			if trans == metaTrans {
				return "", RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewATMTransactionCreator(getTrans, getId)(newTrans)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		created, err := service.NewATMTransactionCreator(getTrans, getId)(newTrans)
		AssertNoError(t, err)
		Assert(t, created.Id, id, "returned id")
	})
}

// TODO: think about maybe having a separate db transaction for every HTTP request

func TestDepositFinalizer(t *testing.T) {
	mockDB := NewStubDB()
	dd := values.DepositData{
		TransactionId: RandomString(),
		Received: []core.Money{
			RandomPositiveMoney(), RandomPositiveMoney(),
		},
	}

	t.Run("should forward to the general finalizer", func(t *testing.T) {
		tErr := RandomError()
		generalFinalizer := func(gotDB db.DB, tId string, tType tValues.TransactionType, gotMoney []core.Money) error {
			if gotDB == mockDB && tId == dd.TransactionId && tType == tValues.Deposit && reflect.DeepEqual(dd.Received, gotMoney) {
				return tErr
			}
			panic("unexpected")
		}
		gotErr := service.NewDepositFinalizer(generalFinalizer)(mockDB, dd)
		AssertError(t, gotErr, tErr)
	})
}

func TestWithdrawalFinalizer(t *testing.T) {
	mockDB := NewStubDB()
	wd := values.WithdrawalData{
		TransactionId: RandomString(),
		Money:         RandomNegativeMoney(),
	}
	t.Run("should forward to the general finalizer", func(t *testing.T) {
		tErr := RandomError()
		generalFinalizer := func(gotDB db.DB, tId string, tType tValues.TransactionType, gotMoney []core.Money) error {
			if gotDB == mockDB && tId == wd.TransactionId && tType == tValues.Withdrawal && reflect.DeepEqual(gotMoney, []core.Money{wd.Money}) {
				return tErr
			}
			panic("unexpected")
		}
		gotErr := service.NewWithdrawalFinalizer(generalFinalizer)(mockDB, wd)
		AssertError(t, gotErr, tErr)
	})
}



func TestGeneralFinalizer(t *testing.T) {
	mockDB := NewStubDB()
	tId := RandomString()
	m := []core.Money{
		{Currency: "USD", Amount: core.NewMoneyAmount(42)},
		{Currency: "RUB", Amount: core.NewMoneyAmount(330.33)},
	}
	metaTrans := tValues.MetaTrans{
		Type:   tValues.Deposit,
		UserId: RandomString(),
	}
	tType := metaTrans.Type
	wantT := []tValues.Transaction{
		{
			Source: tValues.TransSource{
				Type:   tValues.Cash,
				Detail: "",
			},
			UserId: metaTrans.UserId,
			Money:  m[0],
		},
		{
			Source: tValues.TransSource{
				Type:   tValues.Cash,
				Detail: "",
			},
			UserId: metaTrans.UserId,
			Money:  m[1],
		},
	}

	validateMetaTrans := func(gotId string, gotType tValues.TransactionType) (tValues.MetaTrans, error) {
		if gotId == tId && gotType == tType {
			return metaTrans, nil
		}
		panic("unexpected")
	}

	t.Run("error case - validating meta trans throws", func(t *testing.T) {
		tErr := RandomError()
		validateMetaTrans := func(string, tValues.TransactionType) (tValues.MetaTrans, error) {
			return tValues.MetaTrans{}, tErr
		}
		err := service.NewGeneralFinalizer(validateMetaTrans, nil)(mockDB, tId, tType, m)
		AssertError(t, err, tErr)

	})

	t.Run("forward case - forward to multifinalizer", func(t *testing.T) {
		tErr := RandomError()
		finalize := func(gotDB db.DB, gotT []tValues.Transaction) error {
			if gotDB == mockDB && reflect.DeepEqual(gotT, wantT) {
				return tErr
			}
			panic("unexpected")
		}

		err := service.NewGeneralFinalizer(validateMetaTrans, finalize)(mockDB, tId, tType, m)
		AssertError(t, err, tErr)
	})
}

