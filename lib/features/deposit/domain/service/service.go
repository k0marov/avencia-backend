package service

import (
	"github.com/k0marov/avencia-backend/lib/core/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/entities"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/values"
	"time"
)

const ExpDuration = time.Minute * 10

const TransactionTypeClaim = "transaction_type"
const DepositTransactionType = "deposit"

const UserIdClaim = "sub"

type CodeGenerator = func(user auth.User) (code string, expiresAt time.Time, err error)
type CodeVerifier = func(string) (entities.UserInfo, error)
type BanknoteChecker = func(transactionCode string, banknote values.Banknote) bool

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(user auth.User) (string, time.Time, error) {
		claims := map[string]any{
			UserIdClaim:          user.Id,
			TransactionTypeClaim: DepositTransactionType,
		}
		expireAt := time.Now().UTC().Add(ExpDuration)
		code, err := issueJWT(claims, expireAt)
		return code, expireAt, err
	}
}

func NewCodeVerifier(verifyJWT jwt.Verifier) CodeVerifier {
	return func(code string) (entities.UserInfo, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return entities.UserInfo{}, client_errors.InvalidJWT
		}
		if data[TransactionTypeClaim] != DepositTransactionType {
			return entities.UserInfo{}, client_errors.InvalidJWT
		}
		userId, ok := data[UserIdClaim].(string)
		if !ok {
			return entities.UserInfo{}, client_errors.InvalidJWT
		}
		return entities.UserInfo{
			Id: userId,
		}, nil
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(transactionCode string, banknote values.Banknote) bool {
		_, err := verifyCode(transactionCode)
		return err == nil
	}
}
