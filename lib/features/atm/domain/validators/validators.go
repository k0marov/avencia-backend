package validators

import (
	"crypto/subtle"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db"
	"github.com/AvenciaLab/avencia-backend/lib/features/atm/domain/values"
	tStore "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/store"
	tValidators "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/validators"
	tValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/store/mappers"
	wService "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/service"
)

type ATMSecretValidator = func(gotAtmSecret []byte) error
type WithdrawalValidator = func(db.TDB, values.WithdrawalData) error

type baseBanknoteValidator = func(db db.TDB, transId string, expType tValues.TransactionType, banknote core.Money) error
type InsertedBanknoteValidator = func(db.TDB, values.InsertedBanknote) error
type DispensedBanknoteValidator = func(db.TDB, values.DispensedBanknote) error

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

func NewInsertedBanknoteValidator(validate baseBanknoteValidator) InsertedBanknoteValidator {
	return func(db db.TDB, ib values.InsertedBanknote) error {
		return validate(db, ib.TransactionId, tValues.Deposit, ib.Banknote)
	}
}

func NewDispensedBanknoteValidator(validate baseBanknoteValidator) DispensedBanknoteValidator {
	return func(db db.TDB, banknote values.DispensedBanknote) error {
		return validate(db, banknote.TransactionId, tValues.Withdrawal, banknote.Banknote)
	}
}

func NewBaseBanknoteValidator(validate MetaTransByIdValidator, getWallet wService.WalletGetter) baseBanknoteValidator {
	return func(db db.TDB, transId string, expType tValues.TransactionType, banknote core.Money) error {
		trans, err := validate(transId, expType)
		if err != nil {
			return err
		}
		wallet, err := getWallet(db, trans.WalletId)
		if err != nil {
			return err
		}
		if banknote.Currency != wallet.Currency {
			return client_errors.InvalidCurrency
		}
		return nil
	}
}

func NewWithdrawalValidator(validateMeta MetaTransByIdValidator, validateTrans tValidators.TransactionValidator) WithdrawalValidator {
	return func(db db.TDB, wd values.WithdrawalData) error {
		metaTrans, err := validateMeta(wd.TransactionId, tValues.Withdrawal)
		if err != nil {
			return err
		}
		t := tValues.Transaction{
			Source: tValues.TransSource{
				Type: tValues.Cash,
			},
			WalletId: metaTrans.WalletId,
			Money:    wd.Money,
		}
		_, err = validateTrans(db, t)
		return err
	}
}
func NewMetaTransByIdValidator(getTransById tStore.TransactionGetter) MetaTransByIdValidator {
	return newMetaTransValidator(getTransById)
}
func NewMetaTransFromCodeValidator(getTransFromCode mappers.CodeParser) MetaTransByCodeValidator {
	return newMetaTransValidator(getTransFromCode)
}

type anyTransactionGetter = func(someIdentifier string) (tValues.MetaTrans, error)
type metaTransValidator = func(someIdentifier string, wantType tValues.TransactionType) (tValues.MetaTrans, error)

func newMetaTransValidator(getTrans anyTransactionGetter) metaTransValidator {
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
