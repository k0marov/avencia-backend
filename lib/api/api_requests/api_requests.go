package apiRequests

import (
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	"github.com/k0marov/avencia-backend/lib/core/helpers/http_helpers"
	atmValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	transValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
	"net/url"
)

func BanknoteDecoder(_ url.Values, request api.BanknoteCheckRequest) (atmValues.Banknote, error) {
	return atmValues.Banknote{
		TransCode: request.TransactionCode,
		Money: core.Money{
			Currency: core.Currency(request.Currency),
			Amount:   core.NewMoneyAmount(request.Amount),
		},
	}, nil
}

func ATMTransactionDecoder(_ url.Values, request api.FinalizeTransactionRequest) (atmValues.ATMTransaction, error) {
	return atmValues.ATMTransaction{
		ATMSecret: []byte(request.ATMSecret),
		Trans: transValues.Transaction{
			Source: transValues.TransSource{
				Type: transValues.TransSourceType(request.Source.Type),
				Detail: request.Source.Detail,
			},
			UserId: request.UserId,
			Money: core.Money{
				Currency: core.Currency(request.Currency),
				Amount:   core.NewMoneyAmount(request.Amount),
			},
		},
	}, nil
}

func TransferDecoder(user auth.User, _ url.Values, req api.TransferRequest) (transferValues.RawTransfer, error) {
	return transferValues.RawTransfer{
		FromId:  user.Id,
		ToEmail: req.RecipientIdentifier,
		Money: core.Money{
			Currency: core.Currency(req.Currency),
			Amount:   core.NewMoneyAmount(req.Amount),
		},
	}, nil
}

func NewCodeDecoder(user auth.User, query url.Values, _ http_helpers.NoJSONRequest) (atmValues.NewCode, error) {
	transactionType := query.Get(api.TransactionTypeQueryArg)
	if transactionType == "" {
		return atmValues.NewCode{}, client_errors.TransactionTypeNotProvided
	}
	return atmValues.NewCode{
		TransType: atmValues.TransactionType(transactionType),
		User:      user,
	}, nil
}

func CodeForCheckDecoder(query url.Values, req api.CodeRequest) (atmValues.CodeForCheck, error) {
	transactionType := query.Get(api.TransactionTypeQueryArg)
	if transactionType == "" {
		return atmValues.CodeForCheck{}, client_errors.TransactionTypeNotProvided
	}
	return atmValues.CodeForCheck{
		Code:      req.TransactionCode,
		TransType: atmValues.TransactionType(transactionType),
	}, nil
}
