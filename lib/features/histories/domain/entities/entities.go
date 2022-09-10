package entities

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)
type TransEntry struct {
	Source    transValues.TransSource `json:"source"`
	Money     core.Money `json:"money"`
	CreatedAt time.Time `json:"created_at"`
}

