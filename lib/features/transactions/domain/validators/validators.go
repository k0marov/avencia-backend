package validators

import (
	"math"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	limitsService "github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/service"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	walletService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

type TransactionValidator = func(db db.DB, t values.Transaction) (curBalance core.MoneyAmount, err error)

func NewTransactionValidator(checkLimits limitsService.LimitChecker, getBalance walletService.BalanceGetter) TransactionValidator {
	return func(db db.DB, t values.Transaction) (curBalance core.MoneyAmount, err error) {
		if err := checkLimits(db,t); err != nil {
			return core.NewMoneyAmount(0), err
		}


		// TODO: move this block to a separate enoughBalanceValidator  

		// ===== 
		bal, err := getBalance(db, t.UserId, t.Money.Currency)
		if err != nil {
			return core.NewMoneyAmount(0), core_err.Rethrow("getting current balance", err)
		}
		if t.Money.Amount.IsNeg() {
			if bal.Num() < math.Abs(t.Money.Amount.Num()) {
				return core.NewMoneyAmount(0), client_errors.InsufficientFunds
			}
		}
		// ===== 
		
		return bal, nil
	}
}
