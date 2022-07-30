package values

import (
	"github.com/k0marov/avencia-backend/lib/core"
)

type Transaction struct {
	Source TransSource
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


