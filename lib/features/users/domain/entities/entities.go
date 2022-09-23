package entities

import (
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type UserInfo struct {
	User authEntities.DetailedUser
	Wallet walletEntities.Wallet
	Limits limits.Limits
}

