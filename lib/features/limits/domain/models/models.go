package models

import (
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
)

type Withdrawn struct {
	Withdrawn core.Money
	UpdatedAt time.Time
}
