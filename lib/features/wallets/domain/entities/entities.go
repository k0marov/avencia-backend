package entities

import "github.com/AvenciaLab/avencia-backend/lib/core"


const WalletAmountKey = "amount"

type WalletVal struct {
  OwnerId string `json:"owner_id"` 
  Currency core.Currency `json:"currency"`
  Amount core.MoneyAmount `json:"amount"`
}

type Wallet struct {
  Id string `json:"id"`
  WalletVal
}
