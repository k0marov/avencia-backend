package validators

import (
	"crypto/subtle"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
)

// TransCodeValidator err can be a ClientError
type TransCodeValidator = func(code string, wantType values.TransactionType) (userId string, err error)

// ATMSecretValidator err can be a ClientError
type ATMSecretValidator = func(gotAtmSecret []byte) error

func NewTransCodeValidator(verifyJWT jwt.Verifier) TransCodeValidator {
	return func(code string, wantType values.TransactionType) (string, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return "", client_errors.InvalidCode
		}
		if data[values.TransactionTypeClaim] != string(wantType) {
			return "", client_errors.InvalidTransactionType
		}
		userId, ok := data[values.UserIdClaim].(string)
		if !ok {
			return "", client_errors.InvalidCode
		}
		return userId, nil
	}
}

func NewATMSecretValidator(trueATMSecret []byte) ATMSecretValidator {
	return func(gotAtmSecret []byte) error {
		if subtle.ConstantTimeCompare(gotAtmSecret, trueATMSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		return nil
	}
}
