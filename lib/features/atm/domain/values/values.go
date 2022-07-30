package values

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

// TransactionType is either Deposit or Withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal                 = "withdrawal"
)

type NewCode struct {
	TransType TransactionType
	User      auth.User
}

type CodeForCheck struct {
	Code      string
	TransType TransactionType
}

type GeneratedCode struct {
	Code      string
	ExpiresAt time.Time
}

type Banknote struct {
	TransCode string
	Money     core.Money
}


type ATMTransaction struct {
	ATMSecret []byte
	Trans     transValues.Transaction
}

const UserIdClaim = "sub"
const TransactionTypeClaim = "transaction_type"
