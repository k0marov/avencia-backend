package validators_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"testing"
)

func TestTransCodeValidator(t *testing.T) {
	tCode := RandomString()
	tType := RandomTransactionType()
	t.Run("error case - token is invalid", func(t *testing.T) {
		jwtVerifier := func(token string) (map[string]any, error) {
			if token == tCode {
				return nil, RandomError()
			}
			panic("unexpected")
		}
		_, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
		AssertError(t, err, client_errors.InvalidCode)
	})
	t.Run("table test for claims", func(t *testing.T) {
		userId := RandomString()
		tCases := []struct {
			name    string
			claims  map[string]any
			wantErr error
		}{
			{
				"incorrect transction_type claim",
				map[string]any{values.UserIdClaim: "4242", values.TransactionTypeClaim: "random"},
				client_errors.InvalidTransactionType,
			},
			{
				"claims have invalid types",
				map[string]any{values.UserIdClaim: 42, values.TransactionTypeClaim: string(tType)},
				client_errors.InvalidCode,
			},
			{
				"happy case",
				map[string]any{values.UserIdClaim: userId, values.TransactionTypeClaim: string(tType)},
				nil,
			},
		}
		for _, tt := range tCases {
			t.Run(tt.name, func(t *testing.T) {
				jwtVerifier := func(token string) (map[string]any, error) {
					return tt.claims, nil
				}

				gotUserId, err := validators.NewTransCodeValidator(jwtVerifier)(tCode, tType)
				AssertError(t, err, tt.wantErr)
				if tt.wantErr == nil {
					Assert(t, gotUserId, userId, "returned users id")
				}
			})
		}
	})
}

func TestATMSecretValidator(t *testing.T) {
	trueATMSecret := RandomSecret()
	validator := validators.NewATMSecretValidator(trueATMSecret)
	cases := []struct {
		got []byte
		res error
	}{
		{trueATMSecret, nil},
		{RandomSecret(), client_errors.InvalidATMSecret},
		{[]byte(""), client_errors.InvalidATMSecret},
		{[]byte("as;dfk"), client_errors.InvalidATMSecret},
	}

	for _, tt := range cases {
		t.Run(string(tt.got), func(t *testing.T) {
			Assert(t, validator(tt.got), tt.res, "validator result result")
		})
	}

}
