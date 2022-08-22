package values

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

// TODO: TransSource should be added here
type MetaTrans struct {
	TransType TransactionType
	UserId    string
}

const UserIdClaim = "sub"
const TransactionTypeClaim = "transaction_type"

type Transaction struct {
	Source TransSource
	UserId string
	Money  core.Money
}

type GeneratedCode struct {
	Code      string
	ExpiresAt time.Time
}

type TransSourceType string

const (
	CreditCard TransSourceType = "credit-card"
	Cash                       = "cash"
	Crypto                     = "crypto"
	Transfer                   = "transfer"
)

type TransSource struct {
	Type   TransSourceType
	Detail string
}
