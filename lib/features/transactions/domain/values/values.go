package values

import "github.com/k0marov/avencia-backend/lib/core"

type Transaction struct {
	UserId string
	Money  core.Money
}
