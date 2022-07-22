package values

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
)

// TransactionType is either Deposit or Withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

const UserIdClaim = "sub"
const TransactionTypeClaim = "transaction_type"

type Banknote struct {
	Money core.Money
}

type Transaction struct {
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

func NewTransactionData(request api.FinalizeTransactionRequest) Transaction {
	return Transaction{
		UserId: request.UserId,
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.MoneyAmount(request.Amount),
		},
	}
}
