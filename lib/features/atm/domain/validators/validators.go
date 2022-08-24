package validators

import (
	"crypto/subtle"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type ATMSecretValidator = func(gotAtmSecret []byte) error
type InsertedBanknoteValidator = func(db.DB, values.InsertedBanknote) error 
type DispensedBanknoteValidator = func(db.DB, values.DispensedBanknote) error 
type WithdrawalValidator = func(db.DB, values.WithdrawalData) error 

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
	return func(db db.DB, ib values.InsertedBanknote) error {

		return nil
	}
}

// TODO: implement some checks  
func NewDispensedBanknoteValidator() DispensedBanknoteValidator {
	return func(db db.DB, dp values.DispensedBanknote) error {
		return nil 
	}
}

func NewWithdrawalValidator(getTrans tService.TransactionGetter, validate tValidators.TransactionValidator) WithdrawalValidator {
	return func(db db.DB, wd values.WithdrawalData) error {
		metaTrans, err := getTrans(wd.TransactionId)
		if err != nil {
			return core_err.Rethrow("getting transaction data from transaction id", err)
		}
		t := tValues.Transaction{
			Source: tValues.TransSource{
				Type:   tValues.Cash,
			},
			UserId: metaTrans.UserId,
			Money:  wd.Money,
		}
		_, err = validate(db, t) 
		return err
	}
}
