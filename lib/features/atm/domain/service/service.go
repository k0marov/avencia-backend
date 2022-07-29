package service

import (
	"github.com/k0marov/avencia-backend/lib/config/configurable"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade/batch"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/validators"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tService "github.com/k0marov/avencia-backend/lib/features/transactions/domain/service"
	userEntities "github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
	userService "github.com/k0marov/avencia-backend/lib/features/users/domain/service"
	"time"
)

type CodeGenerator = func(values.NewCode) (values.GeneratedCode, error)
type CodeVerifier = func(values.CodeForCheck) (userEntities.UserInfo, error)
type BanknoteChecker = func(values.Banknote) error
type ATMTransactionFinalizer = func(values.ATMTransaction) error

func NewCodeGenerator(issueJWT jwt.Issuer) CodeGenerator {
	return func(newCode values.NewCode) (values.GeneratedCode, error) {
		claims := map[string]any{
			values.UserIdClaim:          newCode.User.Id,
			values.TransactionTypeClaim: newCode.TransType,
		}
		expireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)
		code, err := issueJWT(claims, expireAt)
		return values.GeneratedCode{
			Code:      code,
			ExpiresAt: expireAt,
		}, err
	}
}

func NewCodeVerifier(validate validators.TransCodeValidator, getInfo userService.UserInfoGetter) CodeVerifier {
	return func(code values.CodeForCheck) (userEntities.UserInfo, error) {
		userId, err := validate(code.Code, code.TransType)
		if err != nil {
			return userEntities.UserInfo{}, err
		}
		return getInfo(userId)
	}
}

func NewBanknoteChecker(verifyCode CodeVerifier) BanknoteChecker {
	return func(banknote values.Banknote) error {
		_, err := verifyCode(values.CodeForCheck{
			Code:      banknote.TransCode,
			TransType: values.Deposit,
		})
		// TODO: more banknote checking
		return err
	}
}

func NewATMTransactionFinalizer(validateSecret validators.ATMSecretValidator, runBatch batch.WriteRunner, finalize tService.TransactionFinalizer) ATMTransactionFinalizer {
	return func(atmTrans values.ATMTransaction) error {
		err := validateSecret(atmTrans.ATMSecret)
		if err != nil {
			return err
		}
		return runBatch(func(u fs_facade.BatchUpdater) error {
			return finalize(u, atmTrans.Trans)
		})
	}
}
