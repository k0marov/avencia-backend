package service

import (
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/store"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	userService "github.com/k0marov/avencia-backend/lib/features/user/domain/service"
	"time"
)

const ExpDuration = time.Minute * 10

type CodeGenerator = func(auth.User, values.TransactionType) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string, values.TransactionType) (userEntities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) error
type TransactionFinalizer = func(atmSecret []byte, t values.Transaction) error

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

func NewTransactionFinalizer(validate validators.TransactionValidator, perform store.TransactionPerformer) TransactionFinalizer {
	return func(gotAtmSecret []byte, t values.Transaction) error {
		bal, err := validate(gotAtmSecret, t)
		if err != nil {
			return err
		}
		return perform(bal, t)
	}
}
