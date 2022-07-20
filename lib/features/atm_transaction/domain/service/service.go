package service

import (
	"crypto/subtle"
	"fmt"
	"github.com/k0marov/avencia-backend/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"math"
	"time"
)

const ExpDuration = time.Minute * 10

const TransactionTypeClaimKey = "transaction_type"

// TransactionType is either Deposit or Withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

const UserIdClaim = "sub"

type CodeGenerator = func(auth.User, TransactionType) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string, TransactionType) (entities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) error
type TransactionFinalizer = func(values.TransactionData) error
type TransactionPerformer = func(values.TransactionData) error // yet to be implemented

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(user auth.User, tType TransactionType) (string, time.Time, error) {
		claims := map[string]any{
			UserIdClaim:             user.Id,
			TransactionTypeClaimKey: tType,
		}
		expireAt := time.Now().UTC().Add(ExpDuration)
		code, err := issueJWT(claims, expireAt)
		return code, expireAt, err
	}
}

func NewCodeVerifier(verifyJWT jwt.Verifier) CodeVerifier {
	return func(code string, tType TransactionType) (entities.UserInfo, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return entities.UserInfo{}, client_errors.InvalidCode
		}
		if data[TransactionTypeClaimKey] != tType {
			return entities.UserInfo{}, client_errors.InvalidCode
		}
		userId, ok := data[UserIdClaim].(string)
		if !ok {
			return entities.UserInfo{}, client_errors.InvalidCode
		}
		return entities.UserInfo{
			Id: userId,
		}, nil
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(transactionCode string, banknote values.Banknote) error {
		_, err := verifyCode(transactionCode, Deposit)
		// TODO: more banknote checking
		return err
	}
}

func NewTransactionFinalizer(atmSecret []byte, perform TransactionPerformer) TransactionFinalizer {
	return func(transaction values.TransactionData) error {
		if subtle.ConstantTimeCompare(transaction.ATMSecret, atmSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		return perform(transaction)
	}
}

// StoreBalanceGetter Should return 0 if the wallet field for the given currency is null
type StoreBalanceGetter = func(userId string, currency string) (float64, error)
type StoreBalanceUpdater = func(userId, currency string, newValue float64) error

func NewTransactionPerformer(getBalance StoreBalanceGetter, updateBalance StoreBalanceUpdater) TransactionPerformer {
	return func(t values.TransactionData) error {
		balance, err := getBalance(t.UserId, t.Currency)
		if err != nil {
			return fmt.Errorf("getting current balance: %w", err)
		}
		if t.Amount < 0 {
			if balance < math.Abs(t.Amount) {
				return client_errors.InsufficientFunds
			}
		}
		err = updateBalance(t.UserId, t.Currency, balance+t.Amount)
		if err != nil {
			return fmt.Errorf("updaing balance: %w", err)
		}
		return nil
	}
}
