package apiRequests

import (
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	atmValues "github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	currValues "github.com/AvenciaLab/avencia-backend/lib/features/currencies/domain/values"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	wallets "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
	"github.com/go-chi/chi/v5"
)

func CurrenciesDecoder(_ *http.Request, req api.GetExchangeRatesRequest) (currValues.Currencies, error) {
	return req.Currencies, nil
}

func NewTransDecoder(user authEntities.User, _ *http.Request, req api.GenTransCodeRequest) (tValues.MetaTrans, error) {
	return tValues.MetaTrans{
		Type:     tValues.TransactionType(req.TransactionType),
		CallerId: user.Id,
		WalletId: req.WalletId,
	}, nil
}

func TransDecoder(_ *http.Request, req api.OnTransactionCreateRequest) (atmValues.TransFromQRCode, error) {
	return atmValues.TransFromQRCode{
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

func WalletCreationDecoder(user authEntities.User, _ *http.Request, req api.CreateWalletRequest) (wallets.WalletCreationData, error) {
	return wallets.WalletCreationData{
		UserId:   user.Id,
		Currency: core.Currency(req.Currency),
	}, nil
}

func DepositDataDecoder(_ *http.Request, d api.CompleteDepositRequest) (atmValues.DepositData, error) {
	return atmValues.DepositData{
		TransactionId: d.TransactionId,
		Received:      multiMoneyDecoder(d.Receivables),
	}, nil
}

func WithdrawalDataDecoder(_ *http.Request, w api.StartWithdrawalRequest) (atmValues.WithdrawalData, error) {
	return atmValues.WithdrawalData{
		TransactionId: w.TransactionId,
		Money: core.Money{
			Currency: core.Currency(w.Currency),
			// in the business logic, it is assumed that a withdrawal's amount is negative
			// but in the api, it is always positive, so here we must negate the value
			Amount: core.NewMoneyAmount(-w.Amount),
		},
	}, nil
}

func InsertedBanknoteDecoder(_ *http.Request, b api.BanknoteInsertionRequest) (atmValues.InsertedBanknote, error) {
	return atmValues.InsertedBanknote{
		TransactionId: b.TransactionId,
		Banknote: core.Money{
			Currency: core.Currency(b.Banknote.Currency),
			Amount:   core.NewMoneyAmount(float64(b.Banknote.Denomination)),
		},
		Received: multiMoneyDecoder(b.Receivables),
	}, nil
}

func DispensedBanknoteDecoder(_ *http.Request, b api.BanknoteDispensionRequest) (atmValues.DispensedBanknote, error) {
	return atmValues.DispensedBanknote{
		TransactionId: b.TransactionId,
		Banknote: core.Money{
			Currency: core.Currency(b.Currency),
			Amount:   core.NewMoneyAmount(float64(b.BanknoteDenomination)),
		},
		Remaining: core.NewMoneyAmount(b.RemainingAmount),
		Requested: core.NewMoneyAmount(b.RequestedAmount),
	}, nil
}

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

func TransferDecoder(user authEntities.User, _ *http.Request, req api.TransferRequest) (values.RawTransfer, error) {
	return values.RawTransfer{
		FromId:         user.Id,
		ToEmail:        req.RecipientIdentifier,
		SourceWalletId: req.SourceWalletId,
		Amount:         core.NewMoneyAmount(req.Amount),
	}, nil
}
