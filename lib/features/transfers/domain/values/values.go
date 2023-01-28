package values

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type Transfer struct {
	FromId string
	ToId   string
	FromWallet entities.Wallet
	ToWallet entities.Wallet
	Money core.Money
}

type RawTransfer struct {
	FromId  string
	ToEmail string
	SourceWalletId string 
	Amount   core.MoneyAmount
}
