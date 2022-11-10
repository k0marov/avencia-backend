package apiResponses

import (
	"github.com/AvenciaLab/avencia-api-contract/api"
	atmValues "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	histEntities "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	wEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

func CreatedTransactionEncoder(t atmValues.CreatedTransaction) api.OnTransactionCreateResponse {
	return api.OnTransactionCreateResponse{
		TransactionId: t.Id,
		Customer: api.CustomerResponse{
			Id:        t.UserInfo.User.Id,
			Email:     t.UserInfo.User.Email,
			Mobile:    t.UserInfo.User.PhoneNum,
			FirstName: t.UserInfo.User.DisplayName,
		},
	}
}

func TransCodeEncoder(code transValues.GeneratedCode) api.GenTransCodeResponse {
	return api.GenTransCodeResponse{
		TransactionCode: code.Code,
		ExpiresAt:       code.ExpiresAt.UTC().Unix(),
	}
}

func WalletEncoder(w wEntities.Wallet) api.WalletResponse {
	return api.WalletResponse{
		Id:       w.Id, // TODO: add the id to the wallet entity
		OwnerId:  w.OwnerId,
		Currency: string(w.Currency),
		Amount:   w.Amount.Num(),
	}
}

func WalletsEncoder(wallets []wEntities.Wallet) api.WalletsResponse {
	resp := api.WalletsResponse{}
	for _, w := range wallets {
		resp.Wallets = append(resp.Wallets, WalletEncoder(w))
	}
	return resp
}

func IdEncoder(id string) api.IdResponse {
	return api.IdResponse{Id: id}
}

func LimitsEncoder(l limits.Limits) map[string]api.LimitResponse {
	resp := map[string]api.LimitResponse{}
	for curr, limit := range l {
		resp[string(curr)] = api.LimitResponse{
			Withdrawn: limit.Withdrawn.Num(),
			Max:       limit.Max.Num(),
		}
	}
	return resp
}
func HistoryEncoder(entries []histEntities.TransEntry) api.TransactionHistory {
	respEntries := []api.TransEntry{}
	for _, e := range entries {
		respEntries = append(respEntries, api.TransEntry{
			TransactedAt: e.CreatedAt,
			Source: api.TransactionSource{
				Type:   string(e.Source.Type),
				Detail: e.Source.Detail,
			},
			Money: api.Money{
				Currency: string(e.Money.Currency),
				Amount:   e.Money.Amount.Num(),
			},
		})
	}
	return api.TransactionHistory{
		Entries: respEntries,
	}
}
