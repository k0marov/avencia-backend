package entities

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)
type TransEntry struct {
	Source    transValues.TransSource
	Money     core.Money
	CreatedAt time.Time
}

