package entities

import (
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
	Limits limits.Limits
}
