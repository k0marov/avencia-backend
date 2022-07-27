package configurable

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"time"
)

const TransactionExpDuration = time.Minute * 10

var LimitedCurrencies = map[core.Currency]core.MoneyAmount{
	"USD": core.NewMoneyAmount(1000),
	"RUB": core.NewMoneyAmount(50000), // mainly for tests
}

func IsWithdrawLimitRelevant(withdrawnAt time.Time) bool {
	currentYear := time.Now().Year()
	return withdrawnAt.Year() == currentYear
}
