package core

// TODO: refactor the core directory 

type Currency string // string for now

type Money struct {
	Currency Currency
	Amount   MoneyAmount // could be both negative (means withdrawing) and positive (deposit)
}

func (m Money) Neg() Money {
	return Money{
		Currency: m.Currency, 
		Amount: m.Amount.Neg(),
	}
}
