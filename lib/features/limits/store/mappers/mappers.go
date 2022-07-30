package mappers

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
)

type WithdrawEncoder = func(core.Money) map[string]any
type WithdrawsDecoder = func(map[string]any) (map[core.Currency]core.MoneyAmount, error) 

func WithdrawEncoderImpl(withdraw core.Money) map[string]any {
  return map[string]any{
    string(withdraw.Currency): withdraw.Amount.Num(),
  }
}

func WithdrawsDecoderImpl(raw map[string]any) (map[core.Currency]core.MoneyAmount, error) {
  withdraws := map[core.Currency]core.MoneyAmount{} 
  for curr, amountRaw := range raw {
    amount, err := general_helpers.DecodeFloat(amountRaw) 
    if err != nil {
      return map[core.Currency]core.MoneyAmount{}, err
    }
    withdraws[core.Currency(curr)] = core.NewMoneyAmount(amount)
  }
  return withdraws, nil
}
