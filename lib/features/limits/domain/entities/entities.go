package entities

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

type Limits map[core.Currency]values.Limit
