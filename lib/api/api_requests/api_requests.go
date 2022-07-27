package apiRequests

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	atmValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
)

func BanknoteDecoder(request api.BanknoteCheckRequest) atmValues.Banknote {
	return atmValues.Banknote{
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.NewMoneyAmount(request.Amount),
		},
	}
}

func TransactionDecoder(request api.FinalizeTransactionRequest) atmValues.Transaction {
	return atmValues.Transaction{
		UserId: request.UserId,
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.NewMoneyAmount(request.Amount),
		},
	}
}

func TransferDecoder(user auth.User, req api.TransferRequest) transferValues.RawTransfer {
	return transferValues.RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Currency),
			Amount:   core.NewMoneyAmount(req.Amount),
		},
	}
}
