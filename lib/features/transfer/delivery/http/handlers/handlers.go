package handlers

import (
	"encoding/json"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-backend/lib/core/http_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
	"net/http"
)

// TODO: maybe generalize the http handlers' structure, because they all look the same

func NewTransferHandler(transfer service.Transferer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := http_helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		var transfReq api.TransferRequest
		json.NewDecoder(r.Body).Decode(&transfReq)

		err := transfer(values.NewRawTransfer(user, transfReq))
		if err != nil {
			http_helpers.HandleServiceError(w, err)
			return
		}
	}
}
