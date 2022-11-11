package entities

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)
type HistEntry struct {
	Source    transValues.TransSource `json:"source"`
	Money     core.Money `json:"money"`
	CreatedAt int64 `json:"created_at"` // a unix timestamp 
}

// History entries are sorted from newest to oldest
type History []HistEntry
