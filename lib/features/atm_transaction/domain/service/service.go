package service

import (
	"crypto/subtle"
	"github.com/k0marov/avencia-backend/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"log"
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
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) bool
type TransactionFinalizer = func(values.TransactionData) bool

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
	return func(transactionCode string, banknote values.Banknote) bool {
		_, err := verifyCode(transactionCode, Deposit)
		// TODO: more banknote checking
		return err == nil
	}
}

type TransactionPerformer = func(values.TransactionData) error // yet to be implemented

func NewTransactionFinalizer(atmSecret []byte, perform TransactionPerformer) TransactionFinalizer {
	return func(transaction values.TransactionData) bool {
		if subtle.ConstantTimeCompare(transaction.ATMSecret, atmSecret) == 0 {
			log.Printf("transaction rejected: invalid atm secret")
			return false
		}
		err := perform(transaction)
		return err == nil
	}
}
