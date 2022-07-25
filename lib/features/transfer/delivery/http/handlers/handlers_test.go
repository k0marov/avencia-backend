package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/k0marov/avencia-backend/lib/core/http_test_helpers"
	. "github.com/k0marov/avencia-backend/lib/core/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/transfer/delivery/http/handlers"
	"github.com/k0marov/avencia-backend/lib/features/transfer/domain/values"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTransferHandler(t *testing.T) {
	http_test_helpers.BaseTest401(t, handlers.NewTransferHandler(nil))
	transferReq := RandomTransferRequest()
	body, _ := json.Marshal(transferReq)
	user := RandomUser()
	req := http_test_helpers.AddAuthDataToRequest(http_test_helpers.CreateRequest(bytes.NewReader(body)), user)
	t.Run("happy case", func(t *testing.T) {
		response := httptest.NewRecorder()
		transfered := false
		transferer := func(rawTransfer values.RawTransfer) error {
			if rawTransfer == values.NewRawTransfer(user, transferReq) {
				transfered = true
				return nil
			}
			panic("unexpected")
		}
		handlers.NewTransferHandler(transferer)(response, req)
		AssertStatusCode(t, response, http.StatusOK)
		Assert(t, transfered, true, "transfer comleted")
	})
	http_test_helpers.BaseTestServiceErrorHandling(t, func(err error, response *httptest.ResponseRecorder) {
		transferer := func(values.RawTransfer) error {
			return err
		}
		handlers.NewTransferHandler(transferer)(response, req)
	})
}
