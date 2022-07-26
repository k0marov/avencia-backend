package core

import "math/big"

type MoneyAmount struct {
	num big.Float
}

func NewMoneyAmount(num float64) MoneyAmount {
	return MoneyAmount{num: *big.NewFloat(num)}
}

func (a MoneyAmount) Num() float64 {
	num, _ := a.num.Float64()
	return num
}

func (a MoneyAmount) Neg() MoneyAmount {
	return MoneyAmount{*big.NewFloat(0).Neg(&a.num)}
}

func (a MoneyAmount) Add(b MoneyAmount) MoneyAmount {
	return MoneyAmount{num: *big.NewFloat(0).Add(&a.num, &b.num)}
}

func (a MoneyAmount) Sub(b MoneyAmount) MoneyAmount {
	return MoneyAmount{num: *big.NewFloat(0).Sub(&a.num, &b.num)}
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
	return a.num.Cmp(&b.num) == 1
}
func (a MoneyAmount) IsLess(b MoneyAmount) bool {
	return a.num.Cmp(&b.num) == -1
}
func (a MoneyAmount) IsEqual(b MoneyAmount) bool {
	return a.num.Cmp(&b.num) == 0
}
