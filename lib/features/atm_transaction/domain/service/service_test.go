package service_test

import (
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"testing"
	"time"
)

// TODO: since this file becomes kinda big, create a separate feature "transactions" and move all ATM-independent stuff in there
// TODO: pluralize all feature names

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"reflect"
)

func TestCodeGenerator(t *testing.T) {
	tUser := auth.User{Id: RandomId()}
	tType := RandomTransactionType()
	newCode := values.NewCode{
		TransType: tType,
		User:      tUser,
	}

	wantClaims := map[string]any{
		values.UserIdClaim:          tUser.Id,
		values.TransactionTypeClaim: tType,
	}
	wantExpireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)

	t.Run("forward test", func(t *testing.T) {
		token := RandomString()
		err := RandomError()
		issueJwt := func(gotClaims map[string]any, exp time.Time) (string, error) {
			if reflect.DeepEqual(gotClaims, wantClaims) && TimeAlmostEqual(wantExpireAt, exp) {
				return token, err
			}
			panic("unexpected")
		}
		gotCode, gotErr := service.NewCodeGenerator(issueJwt)(newCode)
		Assert(t, TimeAlmostEqual(gotCode.ExpiresAt, wantExpireAt), true, "the expiration time is Now + ExpDuration")
		AssertError(t, gotErr, err)
		Assert(t, gotCode.Code, token, "returned token")

	})
}

func TestCodeVerifier(t *testing.T) {
	tCodeForCheck := values.CodeForCheck{
		Code:      RandomString(),
		TransType: RandomTransactionType(),
	}
	t.Run("error case - validating the code throws, should rethrow it", func(t *testing.T) {
		err := RandomError()
		codeValidator := func(code string, transType values.TransactionType) (string, error) {
			if code == tCodeForCheck.Code && transType == tCodeForCheck.TransType {
				return "", err
			}
			panic("unexpected")
		}
		_, gotErr := service.NewCodeVerifier(codeValidator, nil)(tCodeForCheck)
		AssertError(t, gotErr, err)
	})
	t.Run("happy case - should forward to userInfoGetter", func(t *testing.T) {
		tUserInfo := RandomUserInfo()
		codeValidator := func(string, values.TransactionType) (string, error) {
			return tUserInfo.Id, nil
		}
		tErr := RandomError()
		userInfoGetter := func(user string) (userEntities.UserInfo, error) {
			if user == tUserInfo.Id {
				return tUserInfo, tErr
			}
			panic("unexpected")
		}
		gotUserInfo, err := service.NewCodeVerifier(codeValidator, userInfoGetter)(tCodeForCheck)
		AssertError(t, err, tErr)
		Assert(t, gotUserInfo, tUserInfo, "returned user info")
	})

}

func TestBanknoteChecker(t *testing.T) {
	banknote := RandomBanknote()
	t.Run("error case - jwt checking throws", func(t *testing.T) {
		verificationErr := RandomError()
		verifyCode := func(values.CodeForCheck) (userEntities.UserInfo, error) {
			return userEntities.UserInfo{}, verificationErr
		}
		err := service.NewBanknoteChecker(verifyCode)(banknote)
		AssertError(t, err, verificationErr)
	})
	t.Run("happy case - jwt checking does not throw", func(t *testing.T) {
		verifyCode := func(codeForCheck values.CodeForCheck) (userEntities.UserInfo, error) {
			if codeForCheck.Code == banknote.TransCode && codeForCheck.TransType == values.Deposit {
				return userEntities.UserInfo{}, nil
			}
			panic("unexpected")
		}
		err := service.NewBanknoteChecker(verifyCode)(banknote)
		AssertNoError(t, err)
	})
}

func TestATMTransactionFinalizer(t *testing.T) {
	atmTrans := values.ATMTransaction{
		ATMSecret: RandomSecret(),
		Trans:     RandomTransactionData(),
	}
	// TODO: remove this callback hell
	stubRunBatch := func(f func(firestore_facade.BatchUpdater) error) error {
		return f(func(*firestore.DocumentRef, map[string]any) error {
			return nil
		})
	}
	t.Run("error case - validation throws", func(t *testing.T) {
		err := RandomError()
		validate := func(atmSecret []byte) error {
			if reflect.DeepEqual(atmSecret, atmTrans.ATMSecret) {
				return err
			}
			panic("unexpected")
		}
		gotErr := service.NewATMTransactionFinalizer(validate, nil, nil)(atmTrans)
		AssertError(t, gotErr, err)
	})
	t.Run("forward case - forward to TransactionFinalizer", func(t *testing.T) {
		validate := func([]byte) error {
			return nil
		}
		err := RandomError()
		finalize := func(u firestore_facade.BatchUpdater, gotTrans values.Transaction) error {
			if gotTrans == atmTrans.Trans {
				return err
			}
			panic("unexpected")
		}
		gotErr := service.NewATMTransactionFinalizer(validate, stubRunBatch, finalize)(atmTrans)
		AssertError(t, gotErr, err)
	})
}

func TestTransactionFinalizer(t *testing.T) {
	batchUpd := func(*firestore.DocumentRef, map[string]any) error { return nil }
	transaction := RandomTransactionData()
	t.Run("error case - validation throws", func(t *testing.T) {
		err := RandomError()
		validate := func(t values.Transaction) (core.MoneyAmount, error) {
			if t == transaction {
				return core.NewMoneyAmount(0), err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransactionFinalizer(validate, nil)(batchUpd, transaction)
		AssertError(t, gotErr, err)
	})
	t.Run("forward case - return whatever performTransaction returns", func(t *testing.T) {
		wantErr := RandomError()
		currentBalance := RandomPosMoneyAmount()
		validate := func(values.Transaction) (core.MoneyAmount, error) {
			return currentBalance, nil
		}
		performTransaction := func(u firestore_facade.BatchUpdater, curBal core.MoneyAmount, trans values.Transaction) error {
			if curBal == currentBalance && trans == transaction {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(validate, performTransaction)(batchUpd, transaction)
		AssertError(t, err, wantErr)
	})
}

func TestTransactionPerformer(t *testing.T) {
	batchUpd := func(*firestore.DocumentRef, map[string]any) error { return nil }
	userId := RandomString()
	curr := RandomCurrency()

	curBalance := core.NewMoneyAmount(100)
	t.Run("deposit", func(t *testing.T) {
		depTrans := values.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: curr,
				Amount:   core.NewMoneyAmount(232.5),
			},
		}
		t.Run("should compute and update balance in case of deposit", func(t *testing.T) {
			balanceUpdated := false
			updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) error {
				if user == userId && currency == curr && newBal.IsEqual(core.NewMoneyAmount(332.5)) {
					balanceUpdated = true
					return nil
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(updBal, nil, nil)(batchUpd, curBalance, depTrans)
			AssertNoError(t, err)
			Assert(t, balanceUpdated, true, "balance was updated")
		})
		t.Run("updating balance throws", func(t *testing.T) {
			updBal := func(firestore_facade.Updater, string, core.Currency, core.MoneyAmount) error {
				return RandomError()
			}
			err := service.NewTransactionPerformer(updBal, nil, nil)(batchUpd, curBalance, depTrans)
			AssertSomeError(t, err)
		})
	})
	t.Run("withdrawal", func(t *testing.T) {
		withdrawTrans := values.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: curr,
				Amount:   core.NewMoneyAmount(-42),
			},
		}
		t.Run("should additionally compute and update withdrawn in case of withdrawal", func(t *testing.T) {
			newWithdrawn := RandomPositiveMoney()
			getNewWithdrawn := func(transaction values.Transaction) (core.Money, error) {
				if transaction == withdrawTrans {
					return newWithdrawn, nil
				}
				panic("unexpected")
			}
			withdrawnUpdated := false
			updWithdrawn := func(b firestore_facade.Updater, user string, value core.Money) error {
				if user == userId && value == newWithdrawn {
					withdrawnUpdated = true
					return nil
				}
				panic("unexpected")
			}
			balanceUpdated := false
			updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) error {
				if newBal.IsEqual(core.NewMoneyAmount(58)) {
					balanceUpdated = true
					return nil
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(updBal, getNewWithdrawn, updWithdrawn)(batchUpd, curBalance, withdrawTrans)
			AssertNoError(t, err)
			Assert(t, balanceUpdated, true, "balance was updated")
			Assert(t, withdrawnUpdated, true, "withdrawn was updated")
		})
		t.Run("getting new withdrawn value throws", func(t *testing.T) {
			getNewWithdrawn := func(values.Transaction) (core.Money, error) {
				return core.Money{}, RandomError()
			}
			err := service.NewTransactionPerformer(nil, getNewWithdrawn, nil)(batchUpd, curBalance, withdrawTrans)
			AssertSomeError(t, err)
		})
		t.Run("updating withdrawn throws", func(t *testing.T) {
			getNewWithdrawn := func(values.Transaction) (core.Money, error) {
				return core.Money{}, nil
			}
			updWithdrawn := func(firestore_facade.Updater, string, core.Money) error {
				return RandomError()
			}
			err := service.NewTransactionPerformer(nil, getNewWithdrawn, updWithdrawn)(batchUpd, curBalance, withdrawTrans)
			AssertSomeError(t, err)
		})
	})
}
