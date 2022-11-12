package validators

import (
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
)

// TransferValidator error may be a ClientError
type TransferValidator = func(values.Transfer) error

func NewTransferValidator() TransferValidator {
	return func(t values.Transfer) error {
		if t.Amount.IsNeg() {
			return client_errors.NegativeTransferAmount
		}
		if t.Amount.IsEqual(core.NewMoneyAmount(0)) {
			return client_errors.TransferringZero
		}
		if t.ToId == t.FromId {
			return client_errors.TransferringToYourself
		}
		if t.FromId != t.SourceWallet.OwnerId {
			return client_errors.Unauthorized
		}
		return nil
	}
}
