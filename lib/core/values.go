package core

type Currency string     // string for now
type MoneyAmount float64 // float64 for now

type Money struct {
	Currency Currency
	Amount   MoneyAmount // could be both negative (means withdrawing) and positive (deposit)
}
