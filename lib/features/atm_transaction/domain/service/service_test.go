package service_test

import (
	"github.com/k0marov/avencia-backend/api/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
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

	tClaims := map[string]any{service.UserIdClaim: "4242", service.TransactionTypeClaimKey: tType}
	tUserInfo := entities.UserInfo{Id: "4242"}

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
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("error case - claims are invalid (e.g. user id is not a string)", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: 42, service.TransactionTypeClaimKey: service.Deposit}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("error case - token has an incorrect transaction_type claim", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: "4242", service.TransactionTypeClaimKey: "random"}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})

	t.Run("happy case", func(t *testing.T) {
		gotUserInfo, err := service.NewCodeVerifier(jwtVerifier)(tCode, tType)
		AssertNoError(t, err)
		Assert(t, gotUserInfo, tUserInfo, "returned user info")
	})

}

func TestBanknoteChecker(t *testing.T) {
	code := RandomString()
	banknote := values.Banknote{
		Currency: RandomString(),
		Amount:   RandomInt(),
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
