package entities

import "github.com/k0marov/avencia-backend/lib/core"

type Wallet map[core.Currency]core.MoneyAmount

func (w Wallet) ToResponse() map[string]float64 {
	r := map[string]float64{}
	for curr, a := range w {
		r[string(curr)] = a.Num()
	}
	return r
}
