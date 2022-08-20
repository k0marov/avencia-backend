package validators

import (
	"crypto/subtle"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type ATMSecretValidator = func(gotAtmSecret []byte) error
type InsertedBanknoteValidator = func(values.InsertedBanknote) error 
type DispensedBanknoteValidator = func(values.DispensedBanknote) error 
type WithdrawalValidator = func(values.WithdrawalData) error 

func NewATMSecretValidator(trueATMSecret []byte) ATMSecretValidator {
	return func(gotAtmSecret []byte) error {
		if subtle.ConstantTimeCompare(gotAtmSecret, trueATMSecret) == 0 {
			return client_errors.InvalidATMSecret
		}
		return nil
	}
}

// TODO: implement some checks  
func NewInsertedBanknoteValidator() InsertedBanknoteValidator {
	return func(values.InsertedBanknote) error {
		return nil
	}
}

// TODO: implement some checks  
func NewDispensedBanknoteValidator() DispensedBanknoteValidator {
	return func(values.DispensedBanknote) error {
		return nil 
	}
}

func NewWithdrawalValidator(getTrans tService.InitTransDataGetter, validate tValidators.TransactionValidator) WithdrawalValidator {
	return func(wd values.WithdrawalData) error {
		initTrans, err := getTrans(wd.TransactionId)
		if err != nil {
			return core_err.Rethrow("getting transaction data from transaction id", err)
		}
		t := tValues.Transaction{
			Source: tValues.TransSource{
				Type:   tValues.Cash,
			},
			UserId: initTrans.UserId,
			Money:  wd.Money,
		}
		_, err = validate(t) 
		return err
	}
}
