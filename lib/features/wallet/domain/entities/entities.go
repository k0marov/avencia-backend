package entities

import "github.com/k0marov/avencia-backend/lib/core"

type Wallet map[core.Currency]core.MoneyAmount

func (w Wallet) ToResponse() (r map[string]float64) {
	for k, v := range w {
		r[string(k)] = v.Num()
	}
	return
}
