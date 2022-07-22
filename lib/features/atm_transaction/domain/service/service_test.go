package service_test

import (
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
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
		userInfoGetter := func(user string) (entities.UserInfo, error) {
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

func TestUserInfoGetter(t *testing.T) {
	userId := RandomString()
	t.Run("error case - getting wallet throws", func(t *testing.T) {
		getWallet := func(user string) (walletEntities.Wallet, error) {
			if user == userId {
				return walletEntities.Wallet{}, RandomError()
			}
			panic("unexpected")
		}
		_, err := service.NewUserInfoGetter(getWallet)(userId)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		wallet := walletEntities.Wallet{"RUB": 1000, "USD": 100.5}
		getWallet := func(string) (walletEntities.Wallet, error) {
			return wallet, nil
		}
		gotInfo, err := service.NewUserInfoGetter(getWallet)(userId)
		AssertNoError(t, err)
		Assert(t, gotInfo, entities.UserInfo{
			Id:     userId,
			Wallet: wallet,
		}, "returned user info")
	})
}

func TestBanknoteChecker(t *testing.T) {
	code := RandomString()
	banknote := RandomBanknote()
	t.Run("error case - jwt checking throws", func(t *testing.T) {
		verificationErr := RandomError()
		verifyCode := func(string, values.TransactionType) (entities.UserInfo, error) {
			return entities.UserInfo{}, verificationErr
		}
		err := service.NewBanknoteChecker(verifyCode)(code, banknote)
		AssertError(t, err, verificationErr)
	})
	t.Run("happy case - jwt checking does not throw", func(t *testing.T) {
		verifyCode := func(gotCode string, tType values.TransactionType) (entities.UserInfo, error) {
			if gotCode == code && tType == values.Deposit {
				return entities.UserInfo{}, nil
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
		validate := func(secret []byte, t values.TransactionData) (core.MoneyAmount, error) {
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
		validate := func([]byte, values.TransactionData) (core.MoneyAmount, error) {
			return currentBalance, nil
		}
		performTransaction := func(userId string, currency core.Currency, newBal core.MoneyAmount) error {
			if userId == transaction.UserId && currency == transaction.Money.Currency && newBal == currentBalance+transaction.Money.Amount {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(validate, performTransaction)(atmSecret, transaction)
		AssertError(t, err, wantErr)
	})
}
