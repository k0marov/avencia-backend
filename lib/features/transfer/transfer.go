package transfer

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/batch"
	atmService "github.com/k0marov/avencia-backend/lib/features/atm_transaction/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	"github.com/k0marov/avencia-backend/lib/features/transfer/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"net/http"
)

func NewTransferHandlerImpl(fsClient *firestore.Client, userFromEmail auth.UserFromEmail, transact atmService.TransactionFinalizer) http.HandlerFunc {
	converter := service.NewTransferConverter(userFromEmail)
	transfer := service.NewTransferer(converter, service.NewTransferValidator(), batch.NewWriteRunner(fsClient), transact)

	return handlers.NewTransferHandler(transfer)
}
