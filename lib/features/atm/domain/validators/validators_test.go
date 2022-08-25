package validators_test

import (
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)


func TestATMSecretValidator(t *testing.T) {
	trueATMSecret := RandomSecret()
	validator := validators.NewATMSecretValidator(trueATMSecret)
	cases := []struct {
		got []byte
		res error
	}{
		{trueATMSecret, nil},
		{RandomSecret(), client_errors.InvalidATMSecret},
		{[]byte(""), client_errors.InvalidATMSecret},
		{[]byte("as;dfk"), client_errors.InvalidATMSecret},
	}

	for _, tt := range cases {
		t.Run(string(tt.got), func(t *testing.T) {
			Assert(t, validator(tt.got), tt.res, "validation result")
		})
	}

}



func TestWithdrawalValidator(t *testing.T) {
	wd := values.WithdrawalData{
		TransactionId: RandomString(),
		Money:         RandomNegativeMoney(),
	}
	initTrans := tValues.MetaTrans{
		Type: RandomTransactionType(),
		UserId:    RandomString(),
	}

	wantTrans := tValues.Transaction{
		Source: tValues.TransSource{
			Type: tValues.Cash,
		},
		UserId: initTrans.UserId,
		Money:  wd.Money,
	}
	mockDB := NewStubDB()
	
	transDataGetter := func(string) (tValues.MetaTrans, error) {
		return initTrans, nil
	}
	t.Run("error case - getting trans data throws an error", func(t *testing.T) {
		transDataGetter := func(transId string) (tValues.MetaTrans, error) {
			if transId == wd.TransactionId {
				return tValues.MetaTrans{}, RandomError()
			}
			panic("unexpected") 
		}
		err := validators.NewWithdrawalValidator(transDataGetter, nil)(mockDB, wd) 
		AssertSomeError(t, err)
	})
	
	t.Run("forward case - forward to TransactionValidator", func(t *testing.T) {
		tErr := RandomError()
		transValidator := func(gotDB  db.DB, trans tValues.Transaction) (core.MoneyAmount, error) {
			if gotDB == mockDB && trans == wantTrans {
				return core.NewMoneyAmount(42), tErr
			} 		
			panic("unexpected")
		}
		err := validators.NewWithdrawalValidator(transDataGetter, transValidator)(mockDB, wd)
		AssertError(t, err, tErr)
	})
}



func TestMetaTransValidator(t *testing.T) {
	id := RandomString() 
	metaTrans := tValues.MetaTrans{
		Type:   tValues.Deposit,
		UserId: RandomString(),
	}
	wantType := tValues.Deposit
	getTrans := func(gotId string) (tValues.MetaTrans, error) {
		if gotId == id {
			return metaTrans, nil
		}
		panic("unexpected")
	}

	t.Run("error case - get trans throws", func(t *testing.T) {
    getTrans := func(string) (tValues.MetaTrans, error) {
    	return tValues.MetaTrans{
    		Type: wantType,
    	}, RandomError()
    }
    _, err := validators.NewMetaTransValidator(getTrans)(id, wantType)
    AssertSomeError(t, err)
	})

	t.Run("error case - transaction type is not right", func(t *testing.T) {
		_, err := validators.NewMetaTransValidator(getTrans)(id, tValues.Withdrawal)
		AssertError(t, err, client_errors.InvalidTransactionType)
	})

	t.Run("happy case", func(t *testing.T) {
		gotTrans, err := validators.NewMetaTransValidator(getTrans)(id, wantType)
		AssertNoError(t, err)
		Assert(t, gotTrans, metaTrans, "returned trans")
	})
}



