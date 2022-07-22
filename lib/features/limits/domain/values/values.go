package values

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"time"
)

type Limit struct {
	Withdrawn core.MoneyAmount
	Max       core.MoneyAmount
}

type WithdrawnWithUpdated struct {
	Withdrawn core.MoneyAmount
	UpdatedAt time.Time
}
