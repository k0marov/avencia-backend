package entities

import (
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	hist "github.com/AvenciaLab/avencia-backend/lib/features/histories/domain/entities"
	"github.com/AvenciaLab/avencia-backend/lib/features/limits"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)

type UserInfo struct {
	User    authEntities.DetailedUser
	Wallets []walletEntities.Wallet
	History hist.History
	Limits  limits.Limits
}
