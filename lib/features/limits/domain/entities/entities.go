package entities

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

type Limits map[core.Currency]values.Limit

func (l Limits) ToResponse() map[string]api.LimitResponse {
	resp := map[string]api.LimitResponse{}
	for curr, limit := range l {
		resp[string(curr)] = api.LimitResponse{
			Withdrawn: limit.Withdrawn.Num(),
			Max:       limit.Max.Num(),
		}
	}
	return resp
}
