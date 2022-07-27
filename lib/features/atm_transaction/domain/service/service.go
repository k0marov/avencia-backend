package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	userService "github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	"time"
)

type CodeGenerator = func(values.NewCode) (values.GeneratedCode, error)
type CodeVerifier = func(values.CodeForCheck) (userEntities.UserInfo, error)
type BanknoteChecker = func(values.Banknote) error
type ATMTransactionFinalizer = func(values.ATMTransaction) error

type TransactionFinalizer = func(u firestore_facade.BatchUpdater, t values.Transaction) error
type transactionPerformer = func(u firestore_facade.BatchUpdater, curBalance core.MoneyAmount, t values.Transaction) error

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(newCode values.NewCode) (values.GeneratedCode, error) {
		claims := map[string]any{
			values.UserIdClaim:          newCode.User.Id,
			values.TransactionTypeClaim: newCode.TransType,
		}
		expireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)
		code, err := issueJWT(claims, expireAt)
		return values.GeneratedCode{
			Code:      code,
			ExpiresAt: expireAt,
		}, err
	}
}

func NewCodeVerifier(validate validators.TransCodeValidator, getInfo userService.UserInfoGetter) CodeVerifier {
	return func(code values.CodeForCheck) (userEntities.UserInfo, error) {
		userId, err := validate(code.Code, code.TransType)
		if err != nil {
			return userEntities.UserInfo{}, err
		}
		return getInfo(userId)
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(banknote values.Banknote) error {
		_, err := verifyCode(values.CodeForCheck{
			Code:      banknote.TransCode,
			TransType: values.Deposit,
		})
		// TODO: more banknote checking
		return err
	}
}

func NewATMTransactionFinalizer(validateSecret validators.ATMSecretValidator, runBatch batch.WriteRunner, finalize TransactionFinalizer) ATMTransactionFinalizer {
	return func(atmTrans values.ATMTransaction) error {
		err := validateSecret(atmTrans.ATMSecret)
		if err != nil {
			return err
		}
		return runBatch(func(u firestore_facade.BatchUpdater) error {
			return finalize(u, atmTrans.Trans)
		})
	}
}

func NewTransactionFinalizer(validate validators.TransactionValidator, perform transactionPerformer) TransactionFinalizer {
	return func(u firestore_facade.BatchUpdater, t values.Transaction) error {
		bal, err := validate(t)
		if err != nil {
			return err
		}
		return perform(u, bal, t)
	}
}

// TODO: please simplify this (move updating withdraw to a separate service)
func NewTransactionPerformer(updBal walletStore.BalanceUpdater, getNewWithdrawn limitsService.WithdrawnUpdateGetter, updWithdrawn limitsStore.WithdrawUpdater) transactionPerformer {
	return func(u firestore_facade.BatchUpdater, curBal core.MoneyAmount, t values.Transaction) error {
		if t.Money.Amount.IsNeg() {
			withdrawn, err := getNewWithdrawn(t)
			if err != nil {
				return fmt.Errorf("getting the new 'withdrawn' value: %w", err)
			}
			err = updWithdrawn(u, t.UserId, withdrawn)
			if err != nil {
				return fmt.Errorf("updating withdrawn value: %w", err)
			}
		}
		return updBal(u, t.UserId, t.Money.Currency, curBal.Add(t.Money.Amount))
	}
}
