package transfers

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/firestore_facade/batch"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfers/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/transfers/domain/service"
	"net/http"
)

func NewTransferHandlerImpl(fsClient *firestore.Client, userFromEmail auth.UserFromEmail, transact tService.TransactionFinalizer) http.HandlerFunc {
	converter := service.NewTransferConverter(userFromEmail)
	transfer := service.NewTransferer(converter, service.NewTransferValidator(), batch.NewWriteRunner(fsClient), transact)

	return handlers.NewTransferHandler(transfer)
}
