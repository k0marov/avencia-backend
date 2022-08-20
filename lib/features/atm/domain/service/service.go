package service

import "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"

type TransactionCreator = func(values.NewTransaction) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error 

type DepositFinalizer = func(values.DepositData) error 
type WithdrawalFinalizer = func(values.WithdrawalData) error 











