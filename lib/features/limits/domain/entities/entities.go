package entities

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/values"
)

type Limits map[core.Currency]values.Limit
