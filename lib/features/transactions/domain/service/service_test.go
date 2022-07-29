package service_test

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	"testing"
)

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
		t.Run("should compute and update balance", func(t *testing.T) {
			updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) error {
				if user == userId && currency == curr && newBal.IsEqual(core.NewMoneyAmount(332.5)) {
					return nil
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(updBal, nil)(batchUpd, curBalance, depTrans)
			AssertNoError(t, err)
		})
		t.Run("updating balance throws", func(t *testing.T) {
			updBal := func(firestore_facade.Updater, string, core.Currency, core.MoneyAmount) error {
				return RandomError()
			}
			err := service.NewTransactionPerformer(updBal, nil)(batchUpd, curBalance, depTrans)
			AssertSomeError(t, err)
		})
	})
	t.Run("withdrawal", func(t *testing.T) {
		withdrawTrans := values.Transaction{
			UserId: userId,
			Money: core.Money{
				Currency: curr,
				Amount:   core.NewMoneyAmount(-42.5),
			},
		}
		updBal := func(b firestore_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) error {
			if user == userId && currency == curr && newBal.IsEqual(core.NewMoneyAmount(57.5)) {
				return nil
			}
			panic("unexpected")
		}
		t.Run("updating withdrawn throws", func(t *testing.T) {
			updateWithdrawn := func(_ firestore_facade.Updater, trans values.Transaction) error {
				if trans == withdrawTrans {
					return RandomError()
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(updBal, updateWithdrawn)(batchUpd, curBalance, withdrawTrans)
			AssertSomeError(t, err)
		})
		t.Run("happy case", func(t *testing.T) {
			updateWithdrawn := func(_ firestore_facade.Updater, trans values.Transaction) error {
				return nil
			}
			err := service.NewTransactionPerformer(updBal, updateWithdrawn)(batchUpd, curBalance, withdrawTrans)
			AssertNoError(t, err)
		})
	})
}
