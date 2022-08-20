package validators

import (
	"crypto/subtle"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
)

type ATMSecretValidator = func(gotAtmSecret []byte) error
type InsertedBanknoteValidator = func(values.InsertedBanknote) error 
type DispensedBanknoteValidator = func(values.DispensedBanknote) error 
type WithdrawalValidator = func(values.WithdrawalData) error 

func NewATMSecretValidator(trueATMSecret []byte) ATMSecretValidator {
	return func(gotAtmSecret []byte) error {
		if subtle.ConstantTimeCompare(gotAtmSecret, trueATMSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		return nil
	}
}
