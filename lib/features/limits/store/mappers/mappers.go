package mappers

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
)

type WithdrawEncoder = func(core.MoneyAmount) map[string]any
type WithdrawDecoder = func(fs_facade.Document) (values.WithdrawnModel, error) 
type WithdrawsDecoder = func(fs_facade.Documents) ([]values.WithdrawnModel, error)

func WithdrawEncoderImpl(withdrawn core.MoneyAmount) map[string]any {
  return map[string]any{
    "withdrawn": withdrawn.Num(), 
  }
}

func WithdrawDecoderImpl(doc fs_facade.Document) (values.WithdrawnModel, error) {
  amount, err := general_helpers.DecodeFloat(doc.Data["withdrawn"]) 
  if err != nil {
    return values.WithdrawnModel{}, err
  }
  return values.WithdrawnModel{
  	Withdrawn: core.Money{
  	  Currency: core.Currency(doc.Id),
  	  Amount: core.NewMoneyAmount(amount), 
  	},
  	UpdatedAt: doc.UpdatedAt,
  }, nil
}

func WithdrawsDecoderImpl(docs fs_facade.Documents) ([]values.WithdrawnModel, error) {
  var models []values.WithdrawnModel 
  for _, doc := range docs {
    withdrawn, err := WithdrawDecoderImpl(doc) 
    if err != nil {
      return []values.WithdrawnModel{}, err
    }
    models = append(models, withdrawn)
  }
  return models, nil
}
