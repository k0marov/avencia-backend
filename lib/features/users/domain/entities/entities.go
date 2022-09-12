package entities

import (
	limitsEntities "github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/entities"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
	Limits limitsEntities.Limits
}
