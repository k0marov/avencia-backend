package service

import (
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
)


type TransactionCreator = func(values.NewTrans) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error 

type DepositFinalizer = func(values.DepositData) error 
type WithdrawalFinalizer = func(values.WithdrawalData) error 



// TODO: add validation that user doesn't have any active transaction 
func NewTransactionCreator() TransactionCreator {
	return func(nt values.NewTrans) (values.CreatedTransaction, error) {
		return values.CreatedTransaction{}, nil
	}
}











