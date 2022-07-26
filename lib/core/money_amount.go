package core

import (
	"github.com/shopspring/decimal"
)

type MoneyAmount struct {
	num decimal.Decimal
}

func NewMoneyAmount(num float64) MoneyAmount {
	return MoneyAmount{num: decimal.NewFromFloat(num)}
}

func (a MoneyAmount) Num() float64 {
	num, _ := a.num.Float64()
	return num
}

func (a MoneyAmount) Neg() MoneyAmount {
	return MoneyAmount{a.num.Neg()}
}

func (a MoneyAmount) Add(b MoneyAmount) MoneyAmount {
	return MoneyAmount{num: a.num.Add(b.num)}
}

func (a MoneyAmount) Sub(b MoneyAmount) MoneyAmount {
	return MoneyAmount{num: a.num.Sub(b.num)}
}

func (a MoneyAmount) IsSet() bool {
	zeroValue := MoneyAmount{}
	return !a.IsEqual(zeroValue)
}

// IsPos includes 0
func (a MoneyAmount) IsPos() bool {
	return a.num.Sign() != -1
}
func (a MoneyAmount) IsNeg() bool {
	return !a.IsPos()
}

func (a MoneyAmount) IsBigger(b MoneyAmount) bool {
	return a.num.Cmp(b.num) == 1
}
func (a MoneyAmount) IsLess(b MoneyAmount) bool {
	return a.num.Cmp(b.num) == -1
}
func (a MoneyAmount) IsEqual(b MoneyAmount) bool {
	return a.num.Cmp(b.num) == 0
}
