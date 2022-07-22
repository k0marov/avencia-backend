package service

import (
	"fmt"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	walletService "github.com/k0marov/avencia-backend/lib/features/wallet/domain/service"
	"time"
)

const ExpDuration = time.Minute * 10

type CodeGenerator = func(auth.User, values.TransactionType) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string, values.TransactionType) (entities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) error
type TransactionFinalizer = func(atmSecret []byte, t values.TransactionData) error

// helpers
type transactionPerformer = func(values.TransactionData) error
type userInfoGetter = func(userId string) (entities.UserInfo, error)

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

func NewCodeVerifier(validate validators.TransCodeValidator, getInfo userInfoGetter) CodeVerifier {
	return func(code string, tType values.TransactionType) (entities.UserInfo, error) {
		userId, err := validate(code, tType)
		if err != nil {
			return entities.UserInfo{}, err
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
		_, err := verifyCode(transactionCode, values.Deposit)
		// TODO: more banknote checking
		return err
	}
}

func NewTransactionFinalizer(validate validators.TransactionValidator, perform store.TransactionPerformer) TransactionFinalizer {
	return func(gotAtmSecret []byte, t values.TransactionData) error {
		bal, err := validate(gotAtmSecret, t)
		if err != nil {
			return err
		}
		return perform(t.UserId, t.Money.Currency, bal+t.Money.Amount)
	}
}
