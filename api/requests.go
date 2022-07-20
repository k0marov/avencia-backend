package api

type CodeRequest struct {
	TransactionCode string `json:"transaction_code"`
}
type BanknoteCheckRequest struct {
	TransactionCode string  `json:"transaction_code"`
	Currency        string  `json:"currency"`
	Amount          float64 `json:"amount"`
}
type FinalizeTransactionRequest struct {
	UserId    string  `json:"user_id"`
	ATMSecret string  `json:"atm_secret"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"` // negative value means withdrawal
}
