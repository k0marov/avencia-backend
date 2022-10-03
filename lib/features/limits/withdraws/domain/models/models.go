package models

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
)

type Withdraws = map[core.Currency]WithdrawVal

type WithdrawVal struct{
	Withdrawn core.MoneyAmount `json:"withdrawn"`
	UpdatedAt int64  `json:"updated_at"` // a unix timestamp
}
