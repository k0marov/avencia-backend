package models

import (
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
)

type Withdraws = map[core.Currency]WithdrawVal

type WithdrawVal struct{
	Withdrawn core.MoneyAmount `json:"withdrawn"`
	UpdatedAt time.Time  `json:"updated_at"`
}
