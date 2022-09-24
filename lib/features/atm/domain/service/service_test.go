package service_test

import (
	"reflect"
	"testing"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	uEntities "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
)

func TestATMTransactionCreator(t *testing.T) {
	mockDB := NewStubDB()
	newTrans := values.TransFromQRCode{
		Type:       RandomTransactionType(),
		QRCodeText: RandomString(),
	}
	metaTrans := tValues.MetaTrans{
		Type:   newTrans.Type,
		UserId: RandomString(),
	}
	id := RandomString()
	uInfo := RandomUserInfo()


	validate := func(string, tValues.TransactionType) (tValues.MetaTrans, error) {
		return metaTrans, nil
	}
	t.Run("error case - parsing qr code throws", func(t *testing.T) {
		tErr := RandomError() 
		validate := func(code string, tType tValues.TransactionType) (tValues.MetaTrans, error) {
			if code == newTrans.QRCodeText && tType == newTrans.Type {
				return tValues.MetaTrans{}, tErr
			}
			panic("unexpected")
		}
		_, err := service.NewATMTransactionCreator(validate, nil, nil)(mockDB, newTrans)
		AssertError(t, err, tErr)
	})

	getUser := func(db.TDB, string) (uEntities.UserInfo, error) {
    return uInfo, nil 
	}

	t.Run("error case - getting user info throws", func(t *testing.T) {
		getUser := func(gotDB db.TDB, userId string) (uEntities.UserInfo, error) {
			if gotDB == mockDB && userId ==  metaTrans.UserId {
				return uEntities.UserInfo{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewATMTransactionCreator(validate, getUser, nil)(mockDB, newTrans)
		AssertSomeError(t, err)
	})

	create := func(tValues.MetaTrans) (string, error) {
		return id, nil
	}
	t.Run("error case - creating transaction throws", func(t *testing.T) {
		create := func(trans tValues.MetaTrans) (string, error) {
			if trans == metaTrans {
				return "", RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewATMTransactionCreator(validate, getUser, create)(mockDB, newTrans)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		created, err := service.NewATMTransactionCreator(validate, getUser, create)(mockDB, newTrans)
		AssertNoError(t, err)
		wantCreated := values.CreatedTransaction{
			Id:       id,
			UserInfo: uInfo,
		}
		Assert(t, created, wantCreated, "returned trans info")
	})
}


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
		generalFinalizer := func(gotDB db.TDB, tId string, tType tValues.TransactionType, gotMoney []core.Money) error {
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
		generalFinalizer := func(gotDB db.TDB, tId string, tType tValues.TransactionType, gotMoney []core.Money) error {
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
		finalize := func(gotDB db.TDB, gotT []tValues.Transaction) error {
			if gotDB == mockDB && reflect.DeepEqual(gotT, wantT) {
				return tErr
			}
			panic("unexpected")
		}

		err := service.NewGeneralFinalizer(validateMetaTrans, finalize)(mockDB, tId, tType, m)
		AssertError(t, err, tErr)
	})
}

