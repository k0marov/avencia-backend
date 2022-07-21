package entities

import (
	"github.com/k0marov/avencia-backend/api"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
}

func (u UserInfo) ToResponse() api.UserInfoResponse {
	return api.UserInfoResponse{Id: u.Id, Wallet: u.Wallet}
}
