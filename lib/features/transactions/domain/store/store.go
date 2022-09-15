package store

import "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"

type TransactionCreator = func(trans values.MetaTrans) (id string, err error)
type TransactionGetter = func(transactionId string) (values.MetaTrans, error)
