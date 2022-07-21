package service_test

import (
	"github.com/k0marov/avencia-backend/api/client_errors"
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
		service.UserIdClaim:             tUser.Id,
		service.TransactionTypeClaimKey: tType,
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
	userId := "4242"
	tClaims := map[string]any{service.UserIdClaim: userId, service.TransactionTypeClaimKey: tType}

	jwtVerifier := func(token string) (map[string]any, error) {
		if token == tCode {
			return tClaims, nil
		}
		panic("unexpected")
	}
	t.Run("error case - token is invalid", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return nil, RandomError()
		}
		_, err := service.NewCodeVerifier(jwtVerifier, nil)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("error case - token has an incorrect transaction_type claim", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: "4242", service.TransactionTypeClaimKey: "random"}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier, nil)(tCode, tType)
		AssertError(t, err, client_errors.InvalidTransactionType)
	})
	t.Run("error case - claims are invalid (e.g. user id is not a string)", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: 42, service.TransactionTypeClaimKey: service.Deposit}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier, nil)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})

	t.Run("happy case - should forward to userInfoGetter", func(t *testing.T) {
		tUserInfo := entities.UserInfo{
			Id:     userId,
			Wallet: map[string]float64{"USD": 400.0, "BTC": 0.0001},
		}
		tErr := RandomError()
		userInfoGetter := func(user string) (entities.UserInfo, error) {
			if user == userId {
				return tUserInfo, tErr
			}
			panic("unexpected")
		}
		gotUserInfo, err := service.NewCodeVerifier(jwtVerifier, userInfoGetter)(tCode, tType)
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
	banknote := values.Banknote{
		Currency: RandomString(),
		Amount:   RandomFloat(),
	}
	t.Run("error case - jwt checking throws", func(t *testing.T) {
		verificationErr := RandomError()
		verifyCode := func(string, service.TransactionType) (entities.UserInfo, error) {
			return entities.UserInfo{}, verificationErr
		}
		err := service.NewBanknoteChecker(verifyCode)(code, banknote)
		AssertError(t, err, verificationErr)
	})
	t.Run("happy case - jwt checking does not throw", func(t *testing.T) {
		verifyCode := func(gotCode string, tType service.TransactionType) (entities.UserInfo, error) {
			if gotCode == code && tType == service.Deposit {
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
	atmSecret := transaction.ATMSecret
	t.Run("error case - atm secret is invalid", func(t *testing.T) {
		otherAtmSecret := []byte("asdf")
		err := service.NewTransactionFinalizer(otherAtmSecret, nil)(transaction)
		AssertError(t, err, client_errors.InvalidATMSecret)
	})
	t.Run("forward case - return whatever performTransaction returns", func(t *testing.T) {
		wantErr := RandomError()
		performTransaction := func(trans values.TransactionData) error {
			if reflect.DeepEqual(trans, transaction) {
				return wantErr
			}
			panic("unexpected")
		}
		err := service.NewTransactionFinalizer(atmSecret, performTransaction)(transaction)
		AssertError(t, err, wantErr)
	})
}

func TestTransactionPerformer(t *testing.T) {
	balance := 150.0
	userId := RandomString()
	currency := RandomString()

	balanceGetter := func(user string, curr string) (float64, error) {
		if user == userId && curr == currency {
			return balance, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting current balance throws", func(t *testing.T) {
		balanceGetter := func(string, string) (float64, error) {
			return 0, RandomError()
		}
		err := service.NewTransactionPerformer(balanceGetter, nil)(values.TransactionData{})
		AssertSomeError(t, err)
	})
	t.Run("error case - not enough funds for withdrawal", func(t *testing.T) {
		amount := -1000.0
		err := service.NewTransactionPerformer(balanceGetter, nil)(values.TransactionData{
			UserId:   userId,
			Currency: currency,
			Amount:   amount,
		})
		AssertError(t, err, client_errors.InsufficientFunds)
	})

	t.Run("error case - updating balance throws", func(t *testing.T) {
		balanceUpdater := func(user string, curr string, newBalance float64) error {
			if user == userId && curr == currency {
				return RandomError()
			}
			panic("unexpected")
		}
		err := service.NewTransactionPerformer(balanceGetter, balanceUpdater)(values.TransactionData{UserId: userId, Currency: currency})
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		t.Run("for withdrawing", func(t *testing.T) {
			amount := -100.0
			balanceUpdater := func(user string, curr string, newBalance float64) error {
				if FloatsEqual(newBalance, 50.0) {
					return nil
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(balanceGetter, balanceUpdater)(values.TransactionData{
				UserId:   userId,
				Currency: currency,
				Amount:   amount,
			})
			AssertNoError(t, err)
		})
		t.Run("for depositing", func(t *testing.T) {
			amount := 100.0
			balanceUpdater := func(user string, curr string, newBalance float64) error {
				if FloatsEqual(newBalance, 250.0) {
					return nil
				}
				panic("unexpected")
			}
			err := service.NewTransactionPerformer(balanceGetter, balanceUpdater)(values.TransactionData{
				UserId:   userId,
				Currency: currency,
				Amount:   amount,
			})
			AssertNoError(t, err)
		})
	})
}
