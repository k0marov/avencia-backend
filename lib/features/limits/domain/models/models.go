package models

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
)

type Withdraws = map[core.Currency]WithdrawVal

type WithdrawVal struct{
	Withdrawn core.MoneyAmount 
	UpdatedAt time.Time 
}
