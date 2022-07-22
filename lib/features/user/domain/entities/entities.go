package entities

import (
	"github.com/k0marov/avencia-api-contract/api"
	limitsEntities "github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
	Limits limitsEntities.Limits
}

func (u UserInfo) ToResponse() api.UserInfoResponse {
	return api.UserInfoResponse{
		Id:     u.Id,
		Wallet: u.Wallet.ToResponse(),
		Limits: u.Limits.ToResponse(),
	}
}
