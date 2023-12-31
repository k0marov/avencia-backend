package values

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	userEntities "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
)

type TransFromQRCode struct {
	Type       tValues.TransactionType
	QRCodeText string
}

type CreatedTransaction struct {
	Id       string
	UserInfo userEntities.UserInfo
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
