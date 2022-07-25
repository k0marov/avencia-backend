package handlers

import (
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"net/http"
)

// TODO: maybe generalize the http handlers' structure, because they all look the same

func NewTransferHandler(transfer service.Transferer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
