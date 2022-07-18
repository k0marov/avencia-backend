package service_test

import (
	"github.com/k0marov/avencia-backend/lib/core/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/service"
	"reflect"
	"testing"
	"time"
)

func TestCodeGenerator(t *testing.T) {
	tUser := auth.User{Id: RandomId()}
	wantClaims := map[string]any{
		service.UserIdClaim:          tUser.Id,
		service.TransactionTypeClaim: service.DepositTransactionType,
	}

	t.Run("forward test", func(t *testing.T) {
		token := RandomString()
		err := RandomError()
		issueJwt := func(gotClaims map[string]any, expDuration time.Duration) (string, error) {
			if reflect.DeepEqual(gotClaims, wantClaims) && expDuration == service.ExpDuration {
				return token, err
			}
			panic("unexpected")
		}
		gotToken, gotErr := service.NewCodeGenerator(issueJwt)(tUser)
		AssertError(t, gotErr, err)
		Assert(t, gotToken, token, "returned token")

	})
}

func TestCodeVerifier(t *testing.T) {
	tCode := RandomString()
	tClaims := map[string]any{service.UserIdClaim: "4242", service.TransactionTypeClaim: service.DepositTransactionType}
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
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode)
		AssertError(t, err, client_errors.InvalidJWT)
	})
	t.Run("error case - token does not contain the needed claims", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: 42, service.TransactionTypeClaim: service.DepositTransactionType}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode)
		AssertError(t, err, client_errors.InvalidJWT)
	})
	t.Run("error case - token has an incorrect transaction_type claim", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{service.UserIdClaim: "4242", service.TransactionTypeClaim: "random"}, nil
		}
		_, err := service.NewCodeVerifier(jwtVerifier)(tCode)
		AssertError(t, err, client_errors.InvalidJWT)
	})

	t.Run("happy case", func(t *testing.T) {
		gotUserInfo, err := service.NewCodeVerifier(jwtVerifier)(tCode)
		AssertNoError(t, err)
		Assert(t, gotUserInfo, tUserInfo, "returned user info")
	})

}
