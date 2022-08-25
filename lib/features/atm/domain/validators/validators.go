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
type MetaTransValidator = func(transId string, wantType tValues.TransactionType) (tValues.MetaTrans, error)  

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

func NewWithdrawalValidator(validateMeta MetaTransValidator, validateTrans tValidators.TransactionValidator) WithdrawalValidator {
	return func(db db.DB, wd values.WithdrawalData) error {
		metaTrans, err := validateMeta(wd.TransactionId, tValues.Withdrawal)
		if err != nil {
			return err
		}
		t := tValues.Transaction{
			Source: tValues.TransSource{
				Type:   tValues.Cash,
			},
			UserId: metaTrans.UserId,
			Money:  wd.Money,
		}
		_, err = validateTrans(db, t) 
		return err
	}
}


func NewMetaTransValidator(getTrans tService.TransactionGetter)  MetaTransValidator {
	return func(transId string, wantType tValues.TransactionType) (tValues.MetaTrans, error) {
		trans, err := getTrans(transId) 
		if err != nil {
			return tValues.MetaTrans{}, core_err.Rethrow("getting trans from trans id", err)
		}
		if trans.Type != wantType {
			return tValues.MetaTrans{}, client_errors.InvalidTransactionType
		}
		return trans, nil 
	}
}



type DeliveryInsertedBanknoteValidator = func(values.InsertedBanknote) error 
type DeliveryDispensedBanknoteValidator = func(values.DispensedBanknote) error 
type DeliveryWithdrawalValidator = func(values.WithdrawalData) error 

func NewDeliveryInsertedBanknoteValidator(db db.DB, validate InsertedBanknoteValidator) DeliveryInsertedBanknoteValidator {
	return func(ib values.InsertedBanknote) error {
		return validate(db, ib)
	}
}
func NewDeliveryDispensedBanknoteValidator(db db.DB, validate DispensedBanknoteValidator) DeliveryDispensedBanknoteValidator {
	return func(b values.DispensedBanknote) error {
		return validate(db, b)
	}
}
func NewDeliveryWithdrawalValidator(db db.DB, validate WithdrawalValidator) DeliveryWithdrawalValidator {
	return func(wd values.WithdrawalData) error {
		return validate(db, wd)
	}
}





