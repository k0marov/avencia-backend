package entities

import (
	limitsEntities "github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
	Limits limitsEntities.Limits
}
