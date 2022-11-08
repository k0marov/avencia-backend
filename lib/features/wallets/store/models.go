package store 

const UserWalletsKey = "wallets"

type UserWalletsModel struct {
  Wallets []string `json:"wallets"`
}
