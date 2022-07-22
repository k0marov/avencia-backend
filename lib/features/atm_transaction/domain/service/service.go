package service

import (
	"crypto/subtle"
	"fmt"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"math"
	"time"
)

const ExpDuration = time.Minute * 10

// TransactionType is either Deposit or Withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

const UserIdClaim = "sub"
const TransactionTypeClaim = "transaction_type"

type CodeGenerator = func(auth.User, TransactionType) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string, TransactionType) (entities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) error
type TransactionFinalizer = func(atmSecret []byte, t values.TransactionData) error

// helpers
type transactionPerformer = func(values.TransactionData) error
type userInfoGetter = func(userId string) (entities.UserInfo, error)

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(user auth.User, tType TransactionType) (string, time.Time, error) {
		claims := map[string]any{
			UserIdClaim:          user.Id,
			TransactionTypeClaim: tType,
		}
		expireAt := time.Now().UTC().Add(ExpDuration)
		code, err := issueJWT(claims, expireAt)
		return code, expireAt, err
	}
}

// TODO: maybe factor out the validation part into a separate validator
func NewCodeVerifier(verifyJWT jwt.Verifier, getInfo userInfoGetter) CodeVerifier {
	return func(code string, tType TransactionType) (entities.UserInfo, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return entities.UserInfo{}, client_errors.InvalidCode
		}
		if data[TransactionTypeClaim] != string(tType) {
			return entities.UserInfo{}, client_errors.InvalidTransactionType
		}
		userId, ok := data[UserIdClaim].(string)
		if !ok {
			return entities.UserInfo{}, client_errors.InvalidCode
		}
		return getInfo(userId)
	}
}

// TODO: add returning the remaining limits
func NewUserInfoGetter(getWallet walletService.WalletGetter) userInfoGetter {
	return func(userId string) (entities.UserInfo, error) {
		wallet, err := getWallet(userId)
		if err != nil {
			return entities.UserInfo{}, fmt.Errorf("getting wallet for user info: %w", err)
		}
		return entities.UserInfo{Id: userId, Wallet: wallet}, nil
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(transactionCode string, banknote values.Banknote) error {
		_, err := verifyCode(transactionCode, Deposit)
		// TODO: more banknote checking
		return err
	}
}

func NewTransactionFinalizer(atmSecret []byte, perform transactionPerformer) TransactionFinalizer {
	return func(gotAtmSecret []byte, t values.TransactionData) error {
		if subtle.ConstantTimeCompare(gotAtmSecret, atmSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		// TODO: add limit check
		return perform(t)
	}
}

func NewTransactionPerformer(getBalance store.BalanceGetter, updateBalance store.BalanceUpdater) transactionPerformer {
	return func(t values.TransactionData) error {
		balance, err := getBalance(t.UserId, t.Money.Currency)
		if err != nil {
			return fmt.Errorf("getting current balance: %w", err)
		}
		if t.Money.Amount < 0 {
			if float64(balance) < math.Abs(float64(t.Money.Amount)) {
				return client_errors.InsufficientFunds
			}
		}
		err = updateBalance(t.UserId, t.Money.Currency, balance+t.Money.Amount)
		if err != nil {
			return fmt.Errorf("updating balance: %w", err)
		}
		return nil
	}
}
