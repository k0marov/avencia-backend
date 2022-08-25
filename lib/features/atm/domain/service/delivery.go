package service

import (
	"github.com/k0marov/avencia-backend/lib/core/db"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
)


// TODO: do something about this repetitiveness

type DeliveryDepositFinalizer = func(values.DepositData) error 
type DeliveryWithdrawalFinalizer = func(values.WithdrawalData) error 
type DeliveryInsertedBanknoteValidator = func(values.InsertedBanknote) error 
type DeliveryDispensedBanknoteValidator = func(values.DispensedBanknote) error 
type DeliveryWithdrawalValidator = func(values.WithdrawalData) error 

func NewDeliveryInsertedBanknoteValidator(db db.DB, validate validators.InsertedBanknoteValidator) DeliveryInsertedBanknoteValidator {
	return func(ib values.InsertedBanknote) error {
		return validate(db, ib)
	}
}
func NewDeliveryDispensedBanknoteValidator(db db.DB, validate validators.DispensedBanknoteValidator) DeliveryDispensedBanknoteValidator {
	return func(b values.DispensedBanknote) error {
		return validate(db, b)
	}
}
func NewDeliveryWithdrawalValidator(db db.DB, validate validators.WithdrawalValidator) DeliveryWithdrawalValidator {
	return func(wd values.WithdrawalData) error {
		return validate(db, wd)
	}
}


func NewDeliveryDepositFinalizer(runT db.TransactionRunner, finalize DepositFinalizer) DeliveryDepositFinalizer {
  return func(dd values.DepositData) error {
    return runT(func(db db.DB) error {
      return finalize(db, dd)
    })
  }
}

func NewDeliveryWithdrawalFinalizer(runT db.TransactionRunner, finalize WithdrawalFinalizer) DeliveryWithdrawalFinalizer {
  return func(wd values.WithdrawalData) error {
    return runT(func(db db.DB) error {
      return finalize(db, wd)
    })
  }
}


