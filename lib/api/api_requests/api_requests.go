package apiRequests

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	atmValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
)

func NewTransDecoder(_ *http.Request, req api.OnTransactionCreateRequest) (atmValues.NewTrans, error) {
	return atmValues.NewTrans{
		Type:       tValues.TransactionType(req.Type),
		QRCodeText: req.QRCodeText,
	}, nil
}



func CancelTransactionDecoder(r *http.Request, _ http_helpers.NoJSONRequest) (transId string, err error) {
	id :=  chi.URLParam(r, "transactionId") 
	if id == "" {
		return "", client_errors.InvalidTransactionId
	}
	return id, nil 
} 

func TransferDecoder(user auth.User, _ *http.Request, req api.TransferRequest) (transferValues.RawTransfer, error) {
	return transferValues.RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Money.Currency),
			Amount:   core.NewMoneyAmount(req.Money.Amount),
		},
	}, nil
}
