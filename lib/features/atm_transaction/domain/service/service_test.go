package service_test

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	"reflect"
	"testing"
	"time"
)

func TestCodeGenerator(t *testing.T) {
	tUser := auth.User{Id: RandomId()}
	tType := RandomTransactionType()

	wantClaims := map[string]any{
		values.UserIdClaim:          tUser.Id,
		values.TransactionTypeClaim: tType,
	}
	wantExpireAt := time.Now().UTC().Add(service.ExpDuration)

	t.Run("forward test", func(t *testing.T) {
		token := RandomString()
		err := RandomError()
		issueJwt := func(gotClaims map[string]any, exp time.Time) (string, error) {
			if reflect.DeepEqual(gotClaims, wantClaims) && TimeAlmostEqual(wantExpireAt, exp) {
				return token, err
			}
			panic("unexpected")
		}
		gotToken, expireAt, gotErr := service.NewCodeGenerator(issueJwt)(tUser, tType)
		Assert(t, TimeAlmostEqual(expireAt, wantExpireAt), true, "the expiration time is Now + ExpDuration")
		AssertError(t, gotErr, err)
		Assert(t, gotToken, token, "returned token")

	})
}

func TestCodeVerifier(t *testing.T) {
	tCode := RandomString()
	tType := RandomTransactionType()
	t.Run("error case - validating the code throws, should rethrow it", func(t *testing.T) {
		err := RandomError()
		codeValidator := func(code string, transType values.TransactionType) (string, error) {
			if code == tCode && transType == tType {
				return "", err
			}
			panic("unexpected")
		}
		_, gotErr := service.NewCodeVerifier(codeValidator, nil)(tCode, tType)
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
		gotUserInfo, err := service.NewCodeVerifier(codeValidator, userInfoGetter)(tCode, tType)
		AssertError(t, err, tErr)
		Assert(t, gotUserInfo, tUserInfo, "returned user info")
	})

}

func TestBanknoteChecker(t *testing.T) {
	code := RandomString()
	banknote := RandomBanknote()
	t.Run("error case - jwt checking throws", func(t *testing.T) {
		verificationErr := RandomError()
		verifyCode := func(string, values.TransactionType) (userEntities.UserInfo, error) {
			return userEntities.UserInfo{}, verificationErr
		}
		err := service.NewBanknoteChecker(verifyCode)(code, banknote)
		AssertError(t, err, verificationErr)
	})
	t.Run("happy case - jwt checking does not throw", func(t *testing.T) {
		verifyCode := func(gotCode string, tType values.TransactionType) (userEntities.UserInfo, error) {
			if gotCode == code && tType == values.Deposit {
				return userEntities.UserInfo{}, nil
			}
			panic("unexpected")
		}
		err := service.NewBanknoteChecker(verifyCode)(code, banknote)
		AssertNoError(t, err)
	})
}

func TestTransactionFinalizer(t *testing.T) {
	transaction := RandomTransactionData()
	atmSecret := RandomSecret()
	t.Run("error case - validation throws", func(t *testing.T) {
		err := RandomError()
		validate := func(secret []byte, t values.Transaction) (core.MoneyAmount, error) {
			if reflect.DeepEqual(secret, atmSecret) && t == transaction {
				return core.MoneyAmount(0), err
			}
			panic("unexpected")
		}
		gotErr := service.NewTransactionFinalizer(validate, nil)(atmSecret, transaction)
		AssertError(t, gotErr, err)
	})
	t.Run("forward case - return whatever performTransaction returns", func(t *testing.T) {
		wantErr := RandomError()
		currentBalance := RandomMoneyAmount()
		validate := func([]byte, values.Transaction) (core.MoneyAmount, error) {
			return currentBalance, nil
		}
		performTransaction := func(curBal core.MoneyAmount, trans values.Transaction) error {
			if curBal == currentBalance && trans == transaction {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(validate, performTransaction)(atmSecret, transaction)
		AssertError(t, err, wantErr)
	})
}

func TestTransactionPerformer(t *testing.T) {
	batchUpdater := func(*firestore.DocumentRef, map[string]any) error { return nil }
	runBatch := func(f func(u firestore_facade.Updater) error) error {
		err := f(batchUpdater)
		return err
	}

	userId := RandomString()
	curr := RandomCurrency()

	curBalance := core.MoneyAmount(100)
	t.Run("should compute and update balance in case of deposit", func(t *testing.T) {
		trans := values.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: curr,
				Amount:   core.MoneyAmount(232.5),
			},
		}
		balanceUpdated := false
		updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) {
			if user == userId && currency == curr && newBal == core.MoneyAmount(332.5) {
				balanceUpdated = true
				return
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(runBatch, updBal, nil, nil)(curBalance, trans)
		AssertNoError(t, err)
		Assert(t, balanceUpdated, true, "balance was updated")
	})
	t.Run("withdrawal", func(t *testing.T) {
		withdrawTrans := values.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: curr,
				Amount:   core.MoneyAmount(-42),
			},
		}
		t.Run("should additionally compute and update withdrawn in case of withdrawal", func(t *testing.T) {
			newWithdrawn := RandomMoney()
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
			updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) {
				if newBal == core.MoneyAmount(58) {
					balanceUpdated = true
					return
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(runBatch, updBal, getNewWithdrawn, updWithdrawn)(curBalance, withdrawTrans)
			AssertNoError(t, err)
			Assert(t, balanceUpdated, true, "balance was updated")
			Assert(t, withdrawnUpdated, true, "withdrawn was updated")
		})
		t.Run("getting new withdrawn value throws", func(t *testing.T) {
			getNewWithdrawn := func(values.Transaction) (core.Money, error) {
				return core.Money{}, RandomError()
			}
			err := service.NewTransactionPerformer(runBatch, nil, getNewWithdrawn, nil)(curBalance, withdrawTrans)
			AssertSomeError(t, err)
		})
		t.Run("updating withdrawn throws", func(t *testing.T) {
			getNewWithdrawn := func(values.Transaction) (core.Money, error) {
				return core.Money{}, RandomError()
			}
			updWithdrawn := func(firestore_facade.Updater, string, core.Money) error {
				return RandomError()
			}
			err := service.NewTransactionPerformer(runBatch, nil, getNewWithdrawn, updWithdrawn)(curBalance, withdrawTrans)
			AssertSomeError(t, err)
		})
	})
	t.Run("should return result of the batchUpdater write", func(t *testing.T) {
		err := RandomError()
		runBatch := func(func(batch firestore_facade.Updater) error) error {
			return err
		}
		gotErr := service.NewTransactionPerformer(runBatch, nil, nil, nil)(RandomMoneyAmount(), RandomTransactionData())
		AssertError(t, gotErr, err)
	})
}
