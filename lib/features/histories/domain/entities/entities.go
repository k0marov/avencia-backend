package entities

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)
type TransEntry struct {
	Source    transValues.TransSource `json:"source"`
	Money     core.Money `json:"money"`
	CreatedAt int64 `json:"created_at"` // a unix timestamp 
}

