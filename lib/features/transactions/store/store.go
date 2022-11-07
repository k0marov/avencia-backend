package store

import (
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/store"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/store/mappers"
)

// TODO: move from using this stupid complicated JWT stuff to a normal db entity 

// The current implementation does not actually create any entities in the database.
// It saves all of the needed data inside a jwt, which is signed and encoded in the id.
// Thanks to this, the API remains stateless + it is more performant.
// But this behavior is subject to change, since it has some limitations.
// For example, is is impossible to invalidate the JWT except for waiting for it to expire,
// Which is currently 10 minutes after it is signed.
// Also:
// TODO: having dots in a transactionId (since it is internally a JWT) may result in having dots in url query params, this may lead to bugs
// Plus, currently the id is just equal to the JWT code.
// But there is a low chance of two JWT's for different transactions being the same.
// So, there is a theoretical possibility of having non-unique ids.
// If this becomes an issue, it is a good idea to just add a random prefix to every id here.


func NewTransactionCreator(genCode mappers.CodeGenerator) store.TransactionCreator {
	return func(trans values.MetaTrans) (id string, err error) {
		code, err := genCode(trans)
		if err != nil {
			return "", core_err.Rethrow("generating code", err)
		}
		return code.Code, nil
	}
}


func NewTransactionGetter(parseCode mappers.CodeParser) store.TransactionGetter {
	return func(transactionId string) (values.MetaTrans, error) {
		code := transactionId // since it internally just the jwt code
		return parseCode(code) 
	}
}
