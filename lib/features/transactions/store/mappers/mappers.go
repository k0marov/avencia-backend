package mappers

import (
	"time"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/jwt"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
)

type CodeGenerator = func(values.MetaTrans) (values.GeneratedCode, error)
type CodeParser = func(code string) (values.MetaTrans, error)

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(trans values.MetaTrans) (values.GeneratedCode, error) {
		claims := map[string]any{
			values.WalletIdClaim:        trans.WalletId,
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

// TODO: somehow simplify this using struct tags and maybe JSON Marshalling
func NewCodeParser(parseJWT jwt.Verifier) CodeParser {
	return func(code string) (values.MetaTrans, error) {
		claims, err := parseJWT(code)
		if err != nil {
			return values.MetaTrans{}, client_errors.InvalidCode
		}

		tType, ok := claims[values.TransactionTypeClaim].(string)
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode
		}
		walletId, ok := claims[values.WalletIdClaim].(string)
		if !ok {
			return values.MetaTrans{}, client_errors.InvalidCode
		}

		return values.MetaTrans{
			Type:   values.TransactionType(tType),
			WalletId: walletId,
		}, nil
	}
}
