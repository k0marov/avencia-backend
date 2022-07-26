package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	limitsService "github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	limitsStore "github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	userService "github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	walletStore "github.com/k0marov/avencia-backend/lib/features/wallet/domain/store"
	"time"
)

const ExpDuration = time.Minute * 10

type CodeGenerator = func(auth.User, values.TransactionType) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string, values.TransactionType) (userEntities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) error

type ATMTransactionFinalizer = func(atmSecret []byte, t values.Transaction) error
type TransactionFinalizer = func(u firestore_facade.BatchUpdater, t values.Transaction) error
type transactionPerformer = func(u firestore_facade.BatchUpdater, curBalance core.MoneyAmount, t values.Transaction) error

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(user auth.User, tType values.TransactionType) (string, time.Time, error) {
		claims := map[string]any{
			values.UserIdClaim:          user.Id,
			values.TransactionTypeClaim: tType,
		}
		expireAt := time.Now().UTC().Add(ExpDuration)
		code, err := issueJWT(claims, expireAt)
		return code, expireAt, err
	}
}

func NewCodeVerifier(validate validators.TransCodeValidator, getInfo userService.UserInfoGetter) CodeVerifier {
	return func(code string, tType values.TransactionType) (userEntities.UserInfo, error) {
		userId, err := validate(code, tType)
		if err != nil {
			return userEntities.UserInfo{}, err
		}
		return getInfo(userId)
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(transactionCode string, banknote values.Banknote) error {
		_, err := verifyCode(transactionCode, values.Deposit)
		// TODO: more banknote checking
		return err
	}
}

func NewATMTransactionFinalizer(validateSecret validators.ATMSecretValidator, runBatch batch.WriteRunner, finalize TransactionFinalizer) ATMTransactionFinalizer {
	return func(atmSecret []byte, t values.Transaction) error {
		err := validateSecret(atmSecret)
		if err != nil {
			return err
		}
		return runBatch(func(u firestore_facade.BatchUpdater) error {
			return finalize(u, t)
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
