package values

type Banknote struct {
	Currency string
	Amount   int
}

type TransactionData struct {
	UserId    string
	ATMSecret []byte
	Amount    int
}
