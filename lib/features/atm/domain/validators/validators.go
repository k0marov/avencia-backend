package validators

import (
	"crypto/subtle"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/transactions/domain/mappers"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	tValidators "github.com/k0marov/avencia-backend/lib/features/transactions/domain/validators"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

type ATMSecretValidator = func(gotAtmSecret []byte) error
type InsertedBanknoteValidator = func(db.DB, values.InsertedBanknote) error 
type DispensedBanknoteValidator = func(db.DB, values.DispensedBanknote) error 
type WithdrawalValidator = func(db.DB, values.WithdrawalData) error 

type MetaTransByIdValidator = func(transId string, wantType tValues.TransactionType) (tValues.MetaTrans, error)  
type MetaTransByCodeValidator = func(code string, wantType tValues.TransactionType) (tValues.MetaTrans, error)

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

func NewWithdrawalValidator(validateMeta MetaTransByIdValidator, validateTrans tValidators.TransactionValidator) WithdrawalValidator {
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

type anyTransactionGetter = func(someIdentifier string) (tValues.MetaTrans, error) 
type metaTransValidator = func(someIdentifier string, wantType tValues.TransactionType) (tValues.MetaTrans, error)
func newMetaTransValidator(getTrans anyTransactionGetter)  metaTransValidator {
	return func(someIdentifier string, wantType tValues.TransactionType) (tValues.MetaTrans, error) {
		trans, err := getTrans(someIdentifier) 
		if err != nil {
			return tValues.MetaTrans{}, core_err.Rethrow("getting trans from an identifier", err)
		}
		if trans.Type != wantType {
			return tValues.MetaTrans{}, client_errors.InvalidTransactionType
		}
		return trans, nil 
	}
}

func NewMetaTransByIdValidator(getTransById tService.TransactionGetter) MetaTransByIdValidator{
	return newMetaTransValidator(getTransById)
}
func NewMetaTransFromCodeValidator(getTransFromCode mappers.CodeParser) MetaTransByCodeValidator {
	return newMetaTransValidator(getTransFromCode)
}



