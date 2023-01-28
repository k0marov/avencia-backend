package validators_test

import (
	"testing"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/validators"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	wService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
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
	metaTrans := tValues.MetaTrans{
		Type:     RandomTransactionType(),
		WalletId: RandomString(),
	}

	wantTrans := tValues.Transaction{
		Source: tValues.TransSource{
			Type:   tValues.Cash,
			Detail: "",
		},
		WalletId: metaTrans.WalletId,
		Money:    wd.Money,
	}
	mockDB := NewStubDB()

	metaTransValidator := func(transactionId string, wantType tValues.TransactionType) (tValues.MetaTrans, error) {
		if transactionId == wd.TransactionId && wantType == tValues.Withdrawal {
			return metaTrans, nil
		}
		panic("unexpected")
	}

	t.Run("error case - validating meta trans throws", func(t *testing.T) {
		tErr := RandomError()
		metaTransValidator := func(string, tValues.TransactionType) (tValues.MetaTrans, error) {
			return tValues.MetaTrans{}, tErr
		}
		err := validators.NewWithdrawalValidator(metaTransValidator, nil)(mockDB, wd)
		AssertError(t, err, tErr)
	})

	t.Run("forward case - forward to TransactionValidator", func(t *testing.T) {
		tErr := RandomError()
		transValidator := func(gotDB db.TDB, trans tValues.Transaction) (core.MoneyAmount, error) {
			if gotDB == mockDB && trans == wantTrans {
				return core.NewMoneyAmount(42), tErr
			}
			panic("unexpected")
		}
		err := validators.NewWithdrawalValidator(metaTransValidator, transValidator)(mockDB, wd)
		AssertError(t, err, tErr)
	})
}

func TestMetaTransValidator(t *testing.T) {
	id := RandomString()
	metaTrans := tValues.MetaTrans{
		Type:     tValues.Deposit,
		WalletId: RandomString(),
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
		_, err := validators.NewMetaTransByIdValidator(getTrans)(id, wantType)
		AssertSomeError(t, err)
	})

	t.Run("error case - transaction type is not right", func(t *testing.T) {
		_, err := validators.NewMetaTransByIdValidator(getTrans)(id, tValues.Withdrawal)
		AssertError(t, err, client_errors.InvalidTransactionType)
	})

	t.Run("happy case", func(t *testing.T) {
		gotTrans, err := validators.NewMetaTransByIdValidator(getTrans)(id, wantType)
		AssertNoError(t, err)
		Assert(t, gotTrans, metaTrans, "returned trans")
	})
}

func TestInsertedBanknoteValidator(t *testing.T) {
	act := func(validate validators.MetaTransByIdValidator, getWallet wService.WalletGetter, tId string, banknote core.Money, db db.TDB) error {
		return validators.NewInsertedBanknoteValidator(validate, getWallet)(db, values.InsertedBanknote{
			TransactionId: tId,
			Banknote:      banknote,
		})
	}
	baseBanknoteValidatorTest(t, tValues.Deposit, act)
}
func TestDispensedBanknoteValidator(t *testing.T) {
	act := func(validate validators.MetaTransByIdValidator, getWallet wService.WalletGetter, tId string, banknote core.Money, db db.TDB) error {
		return validators.NewDispensedBanknoteValidator(validate, getWallet)(db, values.DispensedBanknote{
			TransactionId: tId,
			Banknote:      banknote,
		})
	}
	baseBanknoteValidatorTest(t, tValues.Withdrawal, act)
}

type baseBanknoteValidator = func(validators.MetaTransByIdValidator, wService.WalletGetter, string, core.Money, db.TDB) error

func baseBanknoteValidatorTest(t *testing.T, tType tValues.TransactionType, act baseBanknoteValidator) {
	mockDB := NewStubDB()
	transId := RandomString()
	trans := RandomMetaTrans()
	wallet := RandomWallet()

	t.Run("error case - validating meta trans throws", func(t *testing.T) {
		tErr := RandomError()
		validateTrans := func(gotId string, gotType tValues.TransactionType) (tValues.MetaTrans, error) {
			if gotId == transId && gotType == tType {
				return RandomMetaTrans(), tErr
			}
			panic("unexpected")
		}
		gotErr := act(validateTrans, nil, transId, RandomMoney(), nil)
		AssertError(t, gotErr, tErr)
	})
	validateTrans := func(string, tValues.TransactionType) (tValues.MetaTrans, error) {
		return trans, nil
	}

	t.Run("error case - getting wallet throws", func(t *testing.T) {
		tErr := RandomError()
		getWallet := func(gotDB db.TDB, id string) (wEntities.Wallet, error) {
			if gotDB == mockDB && id == trans.WalletId {
				return RandomWallet(), tErr
			}
			panic("unexpected")
		}
		gotErr := act(validateTrans, getWallet, transId, RandomMoney(), mockDB)
		AssertError(t, gotErr, tErr)
	})

	getWallet := func(db.TDB, string) (wEntities.Wallet, error) {
		return wallet, nil
	}
	t.Run("error case - banknote's currency does not match the wallet currency", func(t *testing.T) {
		banknote := RandomMoney()
		gotErr := act(validateTrans, getWallet, transId, banknote, mockDB)
		AssertError(t, gotErr, client_errors.InvalidCurrency)
	})
	t.Run("happy case", func(t *testing.T) {
  	banknote := core.Money{
  		Currency: wallet.Currency, 
  		Amount: RandomPosMoneyAmount(), 
  	}   
  	err := act(validateTrans, getWallet, transId, banknote, mockDB) 
  	AssertNoError(t, err)
	})
}
