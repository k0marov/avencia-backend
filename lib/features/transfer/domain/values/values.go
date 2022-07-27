package values

import (
	"github.com/k0marov/avencia-backend/lib/core"
)

type Transfer struct {
	FromId string
	ToId   string
	Money  core.Money
}

type RawTransfer struct {
	FromId  string
	ToEmail string
	Money   core.Money
}
