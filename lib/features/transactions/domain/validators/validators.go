package validators

import (
	"math"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallets/domain/service"
)

// TransactionValidator err can be a ClientError
type TransactionValidator = func(t values.Transaction) (curBalance core.MoneyAmount, err error)
// TransCodeValidator err can be a ClientError
type TransCodeValidator = func(code string, wantType values.TransactionType) (values.MetaTrans, error)


func NewTransCodeValidator(verifyJWT jwt.Verifier) TransCodeValidator {
	return func(code string, wantType values.TransactionType) (values.MetaTrans, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return values.MetaTrans{}, client_errors.InvalidCode
		}
		tType, ok := data[values.TransactionTypeClaim].(string) 
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode 
		}
		if tType != string(wantType) {
			return values.MetaTrans{}, client_errors.InvalidTransactionType
		}
		userId, ok := data[values.UserIdClaim].(string)
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode
		}

		trans := values.MetaTrans{
			TransType: values.TransactionType(tType),
			UserId:    userId,
		}

		return trans, nil
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
