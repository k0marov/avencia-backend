package validators

import (
	"math"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	walletService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

type WalletOwnershipValidator = func(db db.TDB, callerId, walletId string) error

type TransactionValidator = func(db db.TDB, t values.Transaction) (curBalance core.MoneyAmount, err error)
type enoughBalanceValidator = func(db db.TDB, t values.Transaction) (curBalance core.MoneyAmount, err error)

func NewWalletOwnershipValidator(getWallet walletService.WalletGetter) WalletOwnershipValidator {
	return func(db db.TDB, callerId, walletId string) error {
		wallet, err := getWallet(db, walletId) 
		if err != nil {
			return core_err.Rethrow("while getting the wallet", err)
		}
		if wallet.OwnerId != callerId {
			return client_errors.Unauthorized
		}
		return nil
	}
}

func NewTransactionValidator(checkLimits limits.LimitChecker, checkBalance enoughBalanceValidator) TransactionValidator {
	return func(db db.TDB, t values.Transaction) (curBalance core.MoneyAmount, err error) {
		if err := checkLimits(db, t); err != nil {
			return core.NewMoneyAmount(0), err
		}

		return checkBalance(db, t)
	}
}

func NewEnoughBalanceValidator(getWallet walletService.WalletGetter) enoughBalanceValidator {
	return func(db db.TDB, t values.Transaction) (curBalance core.MoneyAmount, err error) {
		wallet, err := getWallet(db, t.WalletId)
		if err != nil {
			return core.NewMoneyAmount(0), core_err.Rethrow("getting current balance", err)
		}
		if t.Money.IsNeg() {
			if wallet.Amount.Num() < math.Abs(t.Money.Num()) {
				return core.NewMoneyAmount(0), client_errors.InsufficientFunds
			}
		}
		return wallet.Amount, nil
	}
}
