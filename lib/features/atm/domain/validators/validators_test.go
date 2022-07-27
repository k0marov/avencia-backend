package validators_test

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
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

func TestTransactionValidator(t *testing.T) {
	curBalance := core.NewMoneyAmount(100.0)
	trans := values.Transaction{
		UserId: RandomString(),
		Money: core.Money{
			Currency: RandomCurrency(),
			Amount:   core.NewMoneyAmount(50.0),
		},
	}
	checkLimit := func(t values.Transaction) error {
		return nil
	}
	t.Run("error case - limit checker throws", func(t *testing.T) {
		err := RandomError()
		checkLimit := func(t values.Transaction) error {
			if t == trans {
				return err
			}
			panic("unexpected")
		}
		_, gotErr := validators.NewTransactionValidator(checkLimit, nil)(trans)
		AssertError(t, gotErr, err)
	})
	t.Run("error case - getting balance throws", func(t *testing.T) {
		getBalance := func(string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(0), RandomError()
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertSomeError(t, err)
	})
	t.Run("error case - insufficient funds", func(t *testing.T) {
		getBalance := func(string, core.Currency) (core.MoneyAmount, error) {
			return core.NewMoneyAmount(30.0), nil
		}
		trans := values.Transaction{
			Money: core.Money{
				Amount: core.NewMoneyAmount(-50.0),
			},
		}
		_, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertError(t, err, client_errors.InsufficientFunds)
	})
	t.Run("happy case", func(t *testing.T) {
		getBalance := func(userId string, currency core.Currency) (core.MoneyAmount, error) {
			if userId == trans.UserId && currency == trans.Money.Currency {
				return curBalance, nil
			}
			panic("unexpected")
		}
		bal, err := validators.NewTransactionValidator(checkLimit, getBalance)(trans)
		AssertNoError(t, err)
		Assert(t, bal, curBalance, "returned current balance")
	})
}
