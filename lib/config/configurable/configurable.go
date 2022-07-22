package configurable

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"time"
)

// TODO: move other useful configurables to this package

// These are values that are currently configurable here but may be moved to the config file

var LimitedCurrencies = map[core.Currency]core.MoneyAmount{
	"USD": 1000,
	"RUB": 50000, // mainly for tests
}

func IsWithdrawLimitRelevant(withdrawnAt time.Time) bool {
	currentYear := time.Now().Year()
	return withdrawnAt.After(time.Date(currentYear, 0, 0, 0, 0, 0, 0, time.UTC))
}
