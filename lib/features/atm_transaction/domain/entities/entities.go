package entities

import (
	"github.com/k0marov/avencia-api-contract/api"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
)

type UserInfo struct {
	Id     string
	Wallet walletEntities.Wallet
}

func walletToResponse(w walletEntities.Wallet) (r map[string]float64) {
	for k, v := range w {
		r[string(k)] = float64(v)
	}
	return
}

func (u UserInfo) ToResponse() api.UserInfoResponse {
	return api.UserInfoResponse{Id: u.Id, Wallet: walletToResponse(u.Wallet)}
}
