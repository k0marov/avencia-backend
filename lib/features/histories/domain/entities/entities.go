package entities

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
)

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
