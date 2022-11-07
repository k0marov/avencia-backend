package entities

import "github.com/AvenciaLab/avencia-backend/lib/core"

type Wallet map[core.Currency]core.MoneyAmount 

type WalletInfo struct {
  OwnerId string 
  Money core.Money 
}
