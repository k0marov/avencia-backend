package client_errors

import "fmt"

type ClientError struct {
	DetailCode string `json:"detail_code"`
	HTTPCode   int    `json:"-"`
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("An error which will be displayed to the client: %v %v", ce.HTTPCode, ce.DetailCode)
}

var InvalidCode = ClientError{
	DetailCode: "invalid-code",
	HTTPCode:   400,
}

var InvalidTransactionType = ClientError{
	DetailCode: "invalid-transaction-type",
	HTTPCode:   400,
}

var TransactionTypeNotProvided = ClientError{
	DetailCode: "transaction_type-not-provided",
	HTTPCode:   400,
}

var InvalidATMSecret = ClientError{
	DetailCode: "invalid-atm-secret",
	HTTPCode:   400,
}

var InsufficientFunds = ClientError{
	DetailCode: "insufficient-funds",
	HTTPCode:   400,
}
