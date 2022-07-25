package values

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/auth"
)

type Transfer struct {
	FromId string
	ToId   string
	Money  core.Money
}

type RawTransfer struct {
	FromId  string
	ToEmail string
	Money   core.Money
}

func NewRawTransfer(user auth.User, req api.TransferRequest) RawTransfer {
	return RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Currency),
			Amount:   core.MoneyAmount(req.Amount),
		},
	}
}
