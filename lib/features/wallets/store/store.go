package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/general_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/store"
)

func NewWalletCreator(getDoc db.JsonGetter[UserWalletsModel], updDoc db.JsonUpdater[[]string], setDoc db.JsonSetter[entities.Wallet]) store.WalletCreator {
	return func(db db.TDB, w entities.Wallet) (id string, err error) {
		id = general_helpers.RandomId()
		currWalletsPath := []string{"user_wallet", w.OwnerId}
		currWallets, err := getDoc(db, currWalletsPath)
		if err != nil {
			return "", err
		}
		currWallets.Wallets = append(currWallets.Wallets, id)
		err = updDoc(db, currWalletsPath, UserWalletsKey, currWallets.Wallets)
		if err != nil {
			return "", err
		}
		return id, setDoc(db, []string{"wallets", id}, w)
	}
}

func NewWalletGetter(getDoc db.JsonGetter[entities.Wallet]) store.WalletGetter {
	return func(db db.TDB, walletId string) (entities.Wallet, error) {
		path := []string{"wallets", walletId}
		return getDoc(db, path)
	}
}


// TODO: simplify this
func NewWalletsGetter(getDoc db.JsonGetter[UserWalletsModel], getWallet store.WalletGetter) store.WalletsGetter {
	return func(db db.TDB, userId string) ([]entities.Wallet, error) {
    walletIds, err := getDoc(db, []string{"user_wallets", userId})
    if core_err.IsNotFound(err) {
    	return []entities.Wallet{}, nil
    }
    if err != nil {
    	return []entities.Wallet{}, core_err.Rethrow("while getting the list of user wallet ids", err)
    }
    wallets := []entities.Wallet{}
    for _, id := range walletIds.Wallets {
    	wallet, err := getWallet(db, id)
    	if err != nil {
				return []entities.Wallet{}, core_err.Rethrow("while getting a wallet by id", err)
    	}
    	wallets = append(wallets, wallet)
    }
    return wallets, nil
	}
}


func NewBalanceUpdater(updDoc db.JsonUpdater[core.MoneyAmount]) store.BalanceUpdater {
	return func(db db.TDB, walletId string, newBalance core.MoneyAmount) error {
		path := []string{"wallets", walletId}
		return updDoc(db, path, entities.WalletAmountKey, newBalance)
	}
}
