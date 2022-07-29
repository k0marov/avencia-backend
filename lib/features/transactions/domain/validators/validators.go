package validators

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
	"math"
)

// TransactionValidator err can be a ClientError
type TransactionValidator = func(t values.Transaction) (curBalance core.MoneyAmount, err error)

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
