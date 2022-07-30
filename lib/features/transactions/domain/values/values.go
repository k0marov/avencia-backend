package values

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
)

type Transaction struct {
	UserId string
	Money  core.Money
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

type TransEntry struct {
	Source    TransSource
	Money     core.Money
	CreatedAt time.Time
}


