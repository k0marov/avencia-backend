package service_test

import (
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func TestCodeGenerator(t *testing.T) {
	tUserId := RandomId()
	tType := RandomTransactionType()
	newCode := values.InitTrans{
		TransType: tType,
		UserId:    tUserId,
	}

	wantClaims := map[string]any{
		values.UserIdClaim:          tUserId,
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
		performTransaction := func(u fs_facade.BatchUpdater, curBal core.MoneyAmount, trans values.Transaction) error {
			if curBal == currentBalance && trans == transaction {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(validate, performTransaction)(batchUpd, transaction)
		AssertError(t, err, wantErr)
	})
}

func testTransactionPerfomerForAmount(t *testing.T, transAmount core.MoneyAmount) {
	batchUpd := func(*firestore.DocumentRef, map[string]any) error { return nil }
	curBalance := core.NewMoneyAmount(100)

	trans := values.Transaction{
		Source: RandomTransactionSource(),
		UserId: RandomString(),
		Money: core.Money{
			Currency: RandomCurrency(),
			Amount:   transAmount,
		},
	}

	wantNewBal := curBalance.Add(transAmount)

	var updateWithdrawn limitsService.WithdrawUpdater
	if transAmount.IsNeg() {
		updateWithdrawn = func(fs_facade.Updater, values.Transaction) error {
			return nil
		}
		t.Run("updating withdrawn throws", func(t *testing.T) {
			updateWithdrawn := func(_ fs_facade.Updater, gotTrans values.Transaction) error {
				if gotTrans == trans {
					return RandomError()
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(updateWithdrawn, nil, nil)(batchUpd, curBalance, trans)
			AssertSomeError(t, err)
		})
	}

	addHist := func(u fs_facade.Updater, gotTrans values.Transaction) error {
		if gotTrans == trans {
			return nil
		}
		panic("unexpected")
	}
	t.Run("adding transaction to history throws", func(t *testing.T) {
		addHist := func(fs_facade.Updater, values.Transaction) error {
			return RandomError()
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, nil)(batchUpd, curBalance, trans)
		AssertSomeError(t, err)
	})

	updBal := func(b fs_facade.Updater, user string, currency core.Currency, newBal core.MoneyAmount) error {
		if user == trans.UserId && currency == trans.Money.Currency && newBal.IsEqual(wantNewBal) {
			return nil
		}
		panic("unexpected")
	}
	t.Run("updating balance throws", func(t *testing.T) {
		updBal := func(fs_facade.Updater, string, core.Currency, core.MoneyAmount) error {
			return RandomError()
		}
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(batchUpd, curBalance, trans)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		err := service.NewTransactionPerformer(updateWithdrawn, addHist, updBal)(batchUpd, curBalance, trans)
		AssertNoError(t, err)
	})
}

func TestTransactionPerformer(t *testing.T) {
	t.Run("deposit", func(t *testing.T) {
		testTransactionPerfomerForAmount(t, RandomPosMoneyAmount())
	})
	t.Run("withdrawal", func(t *testing.T) {
		testTransactionPerfomerForAmount(t, RandomNegMoneyAmount())
	})
}
