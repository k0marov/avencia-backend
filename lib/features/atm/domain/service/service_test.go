package service_test

import (
	"testing"

	"github.com/k0marov/avencia-api-contract/api/client_errors"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/service"
	"github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	tValues "github.com/k0marov/avencia-backend/lib/features/transactions/domain/values"
)

func TestATMTransactionCreator(t *testing.T) {
	newTrans := values.NewTrans{
		Type:       RandomTransactionType(),
		QRCodeText: RandomString(),
	}
	metaTrans := tValues.MetaTrans{
		Type:   newTrans.Type,
		UserId: RandomString(),
	}
	id := RandomString()
	getTrans := func(code string) (tValues.MetaTrans, error) {
    return metaTrans, nil
	}
	t.Run("error case - parsing qr code throws", func(t *testing.T) {
	  getTrans := func(code string) (tValues.MetaTrans, error) {
	    if code == newTrans.QRCodeText {
	      return tValues.MetaTrans{}, RandomError()
	    }
	    panic("unexpected")
	  }
	  _, err := service.NewATMTransactionCreator(getTrans, nil)(newTrans)
	  AssertSomeError(t, err)
	})
	t.Run("error case - transaction type is not right", func(t *testing.T) {
	  var wrongMetaTrans tValues.MetaTrans 
	  if newTrans.Type == tValues.Deposit {
	    wrongMetaTrans.Type = tValues.Withdrawal
	  } else {
	    wrongMetaTrans.Type = tValues.Deposit 
	  }
    getTrans := func(string) (tValues.MetaTrans, error) {
      return wrongMetaTrans, nil
    }
    _, err := service.NewATMTransactionCreator(getTrans, nil)(newTrans) 
    AssertError(t, err, client_errors.InvalidTransactionType)
	})

  getId := func(tValues.MetaTrans) (string, error) {
    return id, nil
  }

	t.Run("error case - getting transaction id throws", func(t *testing.T) {
	  getId := func(trans tValues.MetaTrans)(string, error) {
	    if trans == metaTrans {
        return "", RandomError()
	    }
	    panic("unexpected")
	  }
	  _, err := service.NewATMTransactionCreator(getTrans, getId)(newTrans)
	  AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
	  created, err := service.NewATMTransactionCreator(getTrans, getId)(newTrans)
	  AssertNoError(t, err)
	  Assert(t, created.Id, id, "returned id")
	})
}





