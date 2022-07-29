package service_test

import (
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"testing"
	"time"
)

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	userEntities "github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
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
		Assert(t, gotUserInfo, tUserInfo, "returned users info")
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
