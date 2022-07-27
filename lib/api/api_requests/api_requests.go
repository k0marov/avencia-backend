package apiRequests

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	atmValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
	"net/url"
)

func BanknoteDecoder(_ url.Values, request api.BanknoteCheckRequest) (atmValues.Banknote, error) {
	return atmValues.Banknote{
		TransCode: request.TransactionCode,
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.NewMoneyAmount(request.Amount),
		},
	}, nil
}

func ATMTransactionDecoder(_ url.Values, request api.FinalizeTransactionRequest) (atmValues.ATMTransaction, error) {
	return atmValues.ATMTransaction{
		ATMSecret: []byte(request.ATMSecret),
		Trans: atmValues.Transaction{
			UserId: request.UserId,
			Money: core.Money{
				Currency: core.Currency(request.Currency),
				Amount:   core.NewMoneyAmount(request.Amount),
			},
		},
	}, nil
}

func TransferDecoder(user auth.User, _ url.Values, req api.TransferRequest) transferValues.RawTransfer {
	return transferValues.RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Currency),
			Amount:   core.NewMoneyAmount(req.Amount),
		},
	}
}
