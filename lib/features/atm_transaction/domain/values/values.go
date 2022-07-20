package values

import "github.com/k0marov/avencia-backend/api"

type Banknote struct {
	Currency string
	Amount   int
}

type TransactionData struct {
	UserId    string
	ATMSecret []byte
	Currency  string
	Amount    int
}

func NewBanknote(request api.BanknoteCheckRequest) Banknote {
	return Banknote{
		Currency: request.Currency,
		Amount:   request.Amount,
	}
}

func NewTransactionData(request api.FinalizeTransactionRequest) TransactionData {
	return TransactionData{
		UserId:    request.UserId,
		ATMSecret: []byte(request.ATMSecret),
		Currency:  request.Currency,
		Amount:    request.Amount,
	}
}
