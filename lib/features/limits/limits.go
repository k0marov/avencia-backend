package limits

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/limits/domain/store"
	storeImpl "github.com/k0marov/avencia-backend/lib/features/limits/store"
)

type LimitsServices struct {
	GetLimits          service.LimitsGetter
	CheckLimit         service.LimitChecker
	GetWithdrawnUpdate service.WithdrawnUpdateGetter
	UpdateWithdrawn    store.WithdrawUpdater
}

func NewLimitsServicesImpl(fsClient *firestore.Client) LimitsServices {
	storeGetWithdraws := storeImpl.NewWithdrawsGetter(fsClient)
	updateWithdrawn := storeImpl.NewWithdrawUpdater(storeImpl.NewWithdrawsDocGetter(fsClient))

	getLimits := service.NewLimitsGetter(storeGetWithdraws, configurable.LimitedCurrencies)
	return LimitsServices{
		GetLimits:          getLimits,
		CheckLimit:         service.NewLimitChecker(getLimits),
		GetWithdrawnUpdate: service.NewWithdrawnUpdateGetter(getLimits),
		UpdateWithdrawn:    updateWithdrawn,
	}
}
