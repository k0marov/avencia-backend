package service

import (
	"sort"
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/store"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	wallets "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)

type HistoryGetter = func(db db.TDB, userId string) ([]entities.TransEntry, error)
type EntryStorer = func(db.TDB, transValues.Transaction) error


func NewHistoryGetter(getHistory store.HistoryGetter) HistoryGetter {
  return func(db db.TDB, userId string) ([]entities.TransEntry, error) {
  	e, err := getHistory(db, userId) 
  	if err != nil {
  		return []entities.TransEntry{}, core_err.Rethrow("getting history from store", err)
  	}
  	sort.Slice(e, func(i, j int) bool {return e[i].CreatedAt > (e[j].CreatedAt)})
  	return e, nil
  }
}

func NewEntryStorer(getWallet wallets.WalletGetter, storeTrans store.EntryStorer) EntryStorer {
  return func(db db.TDB, t transValues.Transaction) error {
  	wallet, err := getWallet(db, t.WalletId)
		if err != nil {
  		return err
		}
		entry := entities.TransEntry{
			Source:    t.Source,
			Money:     core.Money{
				Currency: wallet.Currency,
				Amount: t.Money,
			},
			CreatedAt: time.Now().Unix(),
		}
		return storeTrans(db, wallet.OwnerId, entry)
  }
}
