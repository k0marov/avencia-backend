package mappers

import (
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/models"
)

type WithdrawEncoder = func(core.MoneyAmount) map[string]any
type WithdrawDecoder = func(db.Document) (models.Withdrawn, error) 
type WithdrawsDecoder = func(db.Documents) ([]models.Withdrawn, error)

func WithdrawEncoderImpl(withdrawn core.MoneyAmount) map[string]any {
  return map[string]any{
    "withdrawn": withdrawn.Num(), 
  }
}

func WithdrawDecoderImpl(doc db.Document) (models.Withdrawn, error) {
  amount, err := general_helpers.DecodeFloat(doc.Data["withdrawn"]) 
  if err != nil {
    return models.Withdrawn{}, err
  }
  return models.Withdrawn{
  	Withdrawn: core.Money{
  	  Currency: core.Currency(doc.Id),
  	  Amount: core.NewMoneyAmount(amount), 
  	},
  	UpdatedAt: doc.UpdatedAt,
  }, nil
}

func WithdrawsDecoderImpl(docs db.Documents) ([]models.Withdrawn, error) {
  var models []models.Withdrawn 
  for _, doc := range docs {
    withdrawn, err := WithdrawDecoderImpl(doc) 
    if err != nil {
      return nil, err
    }
    models = append(models, withdrawn)
  }
  return models, nil
}
