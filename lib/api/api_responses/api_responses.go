package apiResponses

import (
	"github.com/k0marov/avencia-api-contract/api"
	atmValues "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/values"
	limitsEntities "github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	userEntities "github.com/k0marov/avencia-backend/lib/features/user/domain/entities"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallet/domain/entities"
)

func TransCodeEncoder(code atmValues.GeneratedCode) api.CodeResponse {
	return api.CodeResponse{
		TransactionCode: code.Code,
		ExpiresAt:       code.ExpiresAt.UTC().Unix(),
	}
}

func UserInfoEncoder(u userEntities.UserInfo) api.UserInfoResponse {
	return api.UserInfoResponse{
		Id:     u.Id,
		Wallet: WalletEncoder(u.Wallet),
		Limits: LimitsEncoder(u.Limits),
	}
}

func LimitsEncoder(l limitsEntities.Limits) map[string]api.LimitResponse {
	resp := map[string]api.LimitResponse{}
	for curr, limit := range l {
		resp[string(curr)] = api.LimitResponse{
			Withdrawn: limit.Withdrawn.Num(),
			Max:       limit.Max.Num(),
		}
	}
	return resp
}

func WalletEncoder(w walletEntities.Wallet) map[string]float64 {
	r := map[string]float64{}
	for curr, a := range w {
		r[string(curr)] = a.Num()
	}
	return r
}
