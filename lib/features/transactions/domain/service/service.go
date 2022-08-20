package service

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	histService "github.com/k0marov/avencia-backend/lib/features/histories/domain/service"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallets/domain/store"
)

type TransactionFinalizer = func(u fs_facade.BatchUpdater, t values.Transaction) error
type transactionPerformer = func(u fs_facade.BatchUpdater, curBalance core.MoneyAmount, t values.Transaction) error

type CodeGenerator = func(InitTrans) (values.GeneratedCode, error)

type InitTrans struct {
	TransType values.TransactionType
	User      auth.User
}
// TODO: move to transactions
func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(trans InitTrans) (values.GeneratedCode, error) {
		claims := map[string]any{
			values.UserIdClaim:          trans.User.Id,
			values.TransactionTypeClaim: trans.TransType,
		}
		expireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)
		code, err := issueJWT(claims, expireAt)
		return values.GeneratedCode{
			Code:      code,
			ExpiresAt: expireAt,
		}, err
	}
}

func NewTransactionFinalizer(validate validators.TransactionValidator, perform transactionPerformer) TransactionFinalizer {
	return func(u fs_facade.BatchUpdater, t values.Transaction) error {
		bal, err := validate(t)
		if err != nil {
			return err
		}
		return perform(u, bal, t)
	}
}

func NewTransactionPerformer(updateWithdrawn limitsService.WithdrawUpdater, addHist histService.TransStorer, updBal walletStore.BalanceUpdater) transactionPerformer {
	return func(u fs_facade.BatchUpdater, curBal core.MoneyAmount, t values.Transaction) error {
		if t.Money.Amount.IsNeg() {
			err := updateWithdrawn(u, t)
			if err != nil {
				return core_err.Rethrow("updating withdrawn", err)
			}
		}
		err := addHist(u, t) 
		if err != nil {
			return core_err.Rethrow("adding trans to history", err)
		}

		return updBal(u, t.UserId, t.Money.Currency, curBal.Add(t.Money.Amount))
	}
}
