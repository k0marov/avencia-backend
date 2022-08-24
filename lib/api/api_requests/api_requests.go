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

func TransDecoder(_ *http.Request, req api.OnTransactionCreateRequest) (atmValues.NewTrans, error) {
	return atmValues.NewTrans{
		Type:       tValues.TransactionType(req.Type),
		QRCodeText: req.QRCodeText,
	}, nil
}

func CancelTransactionDecoder(r *http.Request, _ http_helpers.NoJSONRequest) (transId string, err error) {
	id := chi.URLParam(r, "transactionId")
	if id == "" {
		return "", client_errors.TransactionIdNotProvided
	}
	return id, nil
}

func WithdrawalDataDecoder(_ *http.Request, w api.StartWithdrawalRequest) (atmValues.WithdrawalData, error) {
	return atmValues.WithdrawalData{
		TransactionId: w.TransactionId,
		Money:         core.Money{
			Currency: core.Currency(w.Currency),
			Amount:   core.NewMoneyAmount(w.Amount),
		},
	}, nil
}

func InsertedBanknoteDecoder(_ *http.Request, b api.BanknoteInsertionRequest) (atmValues.InsertedBanknote, error) {
	return atmValues.InsertedBanknote{
		TransactionId: b.TransactionId,
		Banknote:      core.Money{
			Currency: core.Currency(b.Banknote.Currency), 
			Amount: core.NewMoneyAmount(float64(b.Banknote.Denomination)),
		},
		Received:      multiMoneyDecoder(b.Receivables),
	}, nil
}

func DispensedBanknoteDecoder(_ *http.Request, b api.BanknoteDispensionRequest) (atmValues.DispensedBanknote, error) {
	return atmValues.DispensedBanknote{
		TransactionId: b.TransactionId,
		Banknote:      core.Money{
			Currency: core.Currency(b.Currency),
			Amount:   core.NewMoneyAmount(float64(b.BanknoteDenomination)),
		},
		Remaining:     core.NewMoneyAmount(b.RemainingAmount),
		Requested:     core.NewMoneyAmount(b.RequestedAmount),
	}, nil
}

// TODO: get rid of using type x = func() and replace it with type x func() 


// TODO: move this to some core package 
func moneyDecoder(m api.Money) core.Money {
	return core.Money{
		Currency: core.Currency(m.Currency),
		Amount:   core.NewMoneyAmount(m.Amount),
	}
}

func multiMoneyDecoder(m []api.Money) []core.Money {
  var res []core.Money 
  for _, e := range m {
  	res = append(res, moneyDecoder(e))
  }
  return res
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
