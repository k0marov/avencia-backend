package validators

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

// TransferValidator error may be a ClientError
type TransferValidator = func(values.Transfer) error

func NewTransferValidator() TransferValidator {
	return func(t values.Transfer) error {
		if t.Money.Amount.IsNeg() {
			return client_errors.NegativeTransferAmount
		}
		if t.Money.Amount.IsEqual(core.NewMoneyAmount(0)) {
			return client_errors.TransferingZero
		}
		if t.ToId == t.FromId {
			return client_errors.TransferingToYourself

		}
		return nil
	}
}
