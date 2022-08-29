package mappers

import (
	"fmt"
	"time"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
)

type WithdrawsDecoder = func(db.JsonDocument) (models.Withdraws, error)
type WithdrawEncoder = func(w core.Money, curTime time.Time) map[string]any 

func WithdrawsDecoderImpl(doc db.JsonDocument) (models.Withdraws, error) {
  var withdraws models.Withdraws 
  for k, v := range doc.Data {
  	withdraw, ok := v.(map[string]any)
  	if !ok {
  		return models.Withdraws{}, fmt.Errorf("decoding withdraw val (%v)", v)
  	}
  	amount, err := general_helpers.DecodeFloat(withdraw["amount"])
  	if err != nil {
  		return models.Withdraws{}, core_err.Rethrow("decoding amount", err)
  	}
  	updatedAt, err := general_helpers.DecodeTime(withdraw["updatedAt"])
  	if err != nil {
  		return models.Withdraws{}, core_err.Rethrow("", err)
  	}
    withdraws[core.Currency(k)] = models.WithdrawVal{
    	Withdrawn: core.NewMoneyAmount(amount),
    	UpdatedAt: updatedAt,
    }
  }
  return withdraws, nil
}

func WithdrawEncoderImpl(w core.Money, t time.Time) map[string]any {
  return map[string]any{
  	string(w.Currency): map[string]any{
      "amount": w.Amount.Num(),
      "updatedAt": t.Unix(),
  	},
  }
}
