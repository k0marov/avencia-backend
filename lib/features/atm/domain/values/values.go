package values

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type NewTransaction struct {
	Code string
	Type values.TransactionType
}

type CreatedTransaction struct {
	Id string
	// TODO: add returning user info from the onCreate endpoint
	// UserInfo entities.UserInfo
}

type InsertedBanknote struct {
	TransactionId string
	Banknote      core.Money
	Received      []core.Money
}

type DispensedBanknote struct {
	TransactionId string
	Banknote      core.Money
	Remaining     core.MoneyAmount
	Requested     core.MoneyAmount
}

type DepositData struct {
	TransactionId string
	Received      []core.Money
}

type WithdrawalData struct {
	TransactionId string
	Money         core.Money
}
