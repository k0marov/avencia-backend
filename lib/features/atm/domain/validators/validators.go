package validators

import (
	"crypto/subtle"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
	"math"
)

// TransCodeValidator err can be a ClientError
type TransCodeValidator = func(code string, wantType values.TransactionType) (userId string, err error)

// TransactionValidator err can be a ClientError
type TransactionValidator = func(t values.Transaction) (curBalance core.MoneyAmount, err error)

// ATMSecretValidator err can be a ClientError
type ATMSecretValidator = func(gotAtmSecret []byte) error

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

func NewATMSecretValidator(trueATMSecret []byte) ATMSecretValidator {
	return func(gotAtmSecret []byte) error {
		if subtle.ConstantTimeCompare(gotAtmSecret, trueATMSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		return nil
	}
}

func NewTransactionValidator(checkLimit limitsService.LimitChecker, getBalance walletService.BalanceGetter) TransactionValidator {
	return func(t values.Transaction) (curBalance core.MoneyAmount, err error) {
		if err := checkLimit(t); err != nil {
			return core.NewMoneyAmount(0), err
		}
		bal, err := getBalance(t.UserId, t.Money.Currency)
		if err != nil {
			return core.NewMoneyAmount(0), core_err.Rethrow("getting current balance", err)
		}
		if t.Money.Amount.IsNeg() {
			if bal.Num() < math.Abs(t.Money.Amount.Num()) {
				return core.NewMoneyAmount(0), client_errors.InsufficientFunds
			}
		}
		return bal, nil
	}
}
