package apiRequests

import (
	"net/url"

	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

func TransferDecoder(user auth.User, _ url.Values, req api.TransferRequest) (transferValues.RawTransfer, error) {
	return transferValues.RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Money.Currency),
			Amount:   core.NewMoneyAmount(req.Money.Amount),
		},
	}, nil
}
