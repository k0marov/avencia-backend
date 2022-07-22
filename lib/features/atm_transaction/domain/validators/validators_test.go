package validators_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"testing"
)

func TestTransCodeValidator(t *testing.T) {
	tCode := RandomString()
	tType := RandomTransactionType()
	t.Run("error case - token is invalid", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return nil, RandomError()
		}
		_, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("error case - token has an incorrect transaction_type claim", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{values.UserIdClaim: "4242", values.TransactionTypeClaim: "random"}, nil
		}
		_, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidTransactionType)
	})
	t.Run("error case - claims are invalid (e.g. user id is not a string)", func(t *testing.T) {
		jwtVerifier := func(string) (map[string]any, error) {
			return map[string]any{values.UserIdClaim: 42, values.TransactionTypeClaim: string(tType)}, nil
		}
		_, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("happy case", func(t *testing.T) {
		userId := RandomString()
		tClaims := map[string]any{values.UserIdClaim: userId, values.TransactionTypeClaim: string(tType)}

		jwtVerifier := func(token string) (map[string]any, error) {
			if token == tCode {
				return tClaims, nil
			}
			panic("unexpected")
		}

		gotUserId, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
		AssertNoError(t, err)
		Assert(t, gotUserId, userId, "returned user id")
	})
}
