package values

import (
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
