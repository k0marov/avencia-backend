package core

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type MoneyAmount struct {
	num decimal.Decimal
}

func NewMoneyAmount(num float64) MoneyAmount {
	return MoneyAmount{num: decimal.NewFromFloat(num)}
}

func (m MoneyAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Num())
}
func (m *MoneyAmount) UnmarshalJSON(b []byte) error {
	var num float64
	if err := json.Unmarshal(b, &num); err != nil {
		return err
	}
	m.num = decimal.NewFromFloat(num)
	return nil
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
