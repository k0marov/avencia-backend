package core

// TODO: limit available currencies
type Currency string // string for now

type Money struct {
	Currency Currency
	Amount   MoneyAmount // could be both negative (means withdrawing) and positive (deposit)
}
