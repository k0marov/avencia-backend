package service

import (
	"github.com/k0marov/avencia-backend/lib/core/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/entities"
)

const TransactionTypeClaim = "transaction_type"
const DepositTransactionType = "deposit"

const UserIdClaim = "sub"

type CodeGenerator = func(user auth.User) (string, error)
type CodeVerifier = func(string) (entities.UserInfo, error)

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
