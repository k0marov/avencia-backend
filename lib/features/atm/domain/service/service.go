package service

import (
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tMappers "github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
)


type ATMTransactionCreator = func(values.NewTrans) (values.CreatedTransaction, error)
type TransactionCanceler = func(id string) error 

type DepositFinalizer = func(values.DepositData) error 
type WithdrawalFinalizer = func(values.WithdrawalData) error 



// TODO: add validation that user doesn't have any active transaction 
func NewATMTransactionCreator(getTrans tMappers.CodeParser, getTransId tService.TransactionIdGetter) ATMTransactionCreator {
	return func(nt values.NewTrans) (values.CreatedTransaction, error) {
		trans, err := getTrans(nt.QRCodeText)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting transaction from qr code", err)
		}
		// TODO: maybe move this to a separate validator 
		if trans.Type != nt.Type {
			return values.CreatedTransaction{}, client_errors.InvalidTransactionType
		}
		transId, err := getTransId(trans)
		if err != nil {
			return values.CreatedTransaction{}, core_err.Rethrow("getting the transaction id", err)
		}
		return values.CreatedTransaction{Id: transId}, nil
	}
}











