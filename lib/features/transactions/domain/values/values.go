package values

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

// TODO: TransSource should be added here
type MetaTrans struct {
	Type     TransactionType
	CallerId   string
	WalletId string
}

const WalletIdClaim = "sub"
const CallerIdClaim = "caller_id"
const TransactionTypeClaim = "transaction_type"

type Transaction struct {
	Source   TransSource
	WalletId string
	Money    core.Money
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
