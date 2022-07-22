package validators

import (
	"crypto/subtle"
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"math"
)

// TransCodeValidator err can be a ClientError
type TransCodeValidator = func(code string, wantType values.TransactionType) (userId string, err error)

// TransactionValidator err can be a ClientError
type TransactionValidator = func(gotAtmSecret []byte, t values.TransactionData) (curBalance core.MoneyAmount, err error)

func NewTransCodeValidator(verifyJWT jwt.Verifier) TransCodeValidator {
	return func(code string, wantType values.TransactionType) (string, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return "", client_errors.InvalidCode
		}
		if data[values.TransactionTypeClaim] != string(wantType) {
			return "", client_errors.InvalidTransactionType
		}
		userId, ok := data[values.UserIdClaim].(string)
		if !ok {
			return "", client_errors.InvalidCode
		}
		return userId, nil
	}
}

// TODO: add .toNum method to MoneyAmount

// TODO: add limit check
func NewTransactionValidator(atmSecret []byte, getBalance store.BalanceGetter) TransactionValidator {
	return func(gotSecret []byte, t values.TransactionData) (curBalance core.MoneyAmount, err error) {
		if subtle.ConstantTimeCompare(gotSecret, atmSecret) == 0 {
			return core.MoneyAmount(0), client_errors.InvalidATMSecret
		}
		bal, err := getBalance(t.UserId, t.Money.Currency)
		if err != nil {
			return core.MoneyAmount(0), fmt.Errorf("getting current balance: %w", err)
		}
		if t.Money.Amount < 0 {
			if float64(bal) < math.Abs(float64(t.Money.Amount)) {
				return core.MoneyAmount(0), client_errors.InsufficientFunds
			}
		}
		return bal, nil
	}
}
