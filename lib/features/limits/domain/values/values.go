package values

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"time"
)

type Limit struct {
	Withdrawn core.MoneyAmount
	Max       core.MoneyAmount
}


// TODO: move to models 
type WithdrawnModel struct {
	Withdrawn core.Money
	UpdatedAt time.Time
}
