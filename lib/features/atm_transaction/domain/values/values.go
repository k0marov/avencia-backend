package values

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
)

type Banknote struct {
	Money core.Money
}

type TransactionData struct {
	UserId string
	Money  core.Money
}

func NewBanknote(request api.BanknoteCheckRequest) Banknote {
	return Banknote{
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.MoneyAmount(request.Amount),
		},
	}
}

func NewTransactionData(request api.FinalizeTransactionRequest) TransactionData {
	return TransactionData{
		UserId: request.UserId,
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.MoneyAmount(request.Amount),
		},
	}
}
