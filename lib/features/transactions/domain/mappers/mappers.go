package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type CodeGenerator = func(values.MetaTrans) (values.GeneratedCode, error)
type CodeParser = func(code string) (values.MetaTrans, error) 

type TransIdGenerator = func(code string)  (transId string)
type TransIdParser = func(transactionId string) (code string)

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(trans values.MetaTrans) (values.GeneratedCode, error) {
		claims := map[string]any{
			values.UserIdClaim:          trans.UserId,
			values.TransactionTypeClaim: trans.Type,
		}
		expireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)
		code, err := issueJWT(claims, expireAt)
		return values.GeneratedCode{
			Code:      code,
			ExpiresAt: expireAt,
		}, err
	}
}

func NewCodeParser(parseJWT jwt.Verifier) CodeParser {
	return func(code string) (values.MetaTrans, error) {
		claims, err := parseJWT(code) 
		if err != nil {
			return values.MetaTrans{}, core_err.Rethrow("parsing jwt", err)
		}

		tType, ok := claims[values.TransactionTypeClaim].(string) 
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode 
		}
		userId, ok := claims[values.UserIdClaim].(string)
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode
		}

		return values.MetaTrans{
			Type: values.TransactionType(tType),
			UserId:    userId,
		}, nil
	}
}

func NewTransIdGenerator() TransIdGenerator {
	return func(transCode string) string {
		uuid, _ := uuid.NewUUID()
		return uuid.String() + transCode
	}
}

func NewTransIdParser() TransIdParser {
	return func(transactionId string) string {
		return transactionId[36:]
	}
}
