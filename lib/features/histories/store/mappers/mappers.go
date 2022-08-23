package mappers

import (
	"fmt"

	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/k0marov/avencia-backend/lib/features/histories/domain/entities"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type TransEntryDecoder = func(db.Document) (entities.TransEntry, error) 
type TransEntriesDecoder = func(db.Documents) ([]entities.TransEntry, error) 

type TransEntryEncoder = func(transValues.TransSource, core.Money) map[string]any


func TransEntryEncoderImpl(source transValues.TransSource, money core.Money) map[string]any {
	return map[string]any{
		"source": map[string]any{
			"type": string(source.Type), 
			"detail": source.Detail, 
		},
		"money": map[string]any{
			"currency": string(money.Currency), 
			"amount": money.Amount.Num(),
		},
	}
}

// TODO: replace this nightmare with proper struct tags usage
func TransEntryDecoderImpl(doc db.Document) (entities.TransEntry, error) {
	sourceMap, ok := doc.Data["source"].(map[string]string)
	if !ok {
		return entities.TransEntry{}, fmt.Errorf("decoding transaction source of doc: %+v", doc)
	}
	source := transValues.TransSource{
		Type:   transValues.TransSourceType(sourceMap["type"]),
		Detail: sourceMap["detail"],
	}

	moneyMap, ok := doc.Data["money"].(map[string]any)
	if !ok {
		return entities.TransEntry{}, fmt.Errorf("decoding money of doc: %+v", doc)
	}

	currency, ok := moneyMap["currency"].(string)
	if !ok {
		return entities.TransEntry{}, fmt.Errorf("decoding money currency of doc: %+v", doc)
	}

	amount, err := general_helpers.DecodeFloat(moneyMap["amount"])
	if err != nil {
		return entities.TransEntry{}, core_err.Rethrow("decoding money amount", err)
	}

	money := core.Money{
		Currency: core.Currency(currency),
		Amount:   core.NewMoneyAmount(amount),
	}

	return entities.TransEntry{
		Source:    source,
		Money:     money,
		CreatedAt: doc.CreatedAt,
	}, nil
}


func TransEntriesDecoderImpl(docs db.Documents) ([]entities.TransEntry, error) {
	var entries []entities.TransEntry
	for _, doc := range docs {
		entry, err := TransEntryDecoderImpl(doc) 
		if err != nil {
			return []entities.TransEntry{}, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}


