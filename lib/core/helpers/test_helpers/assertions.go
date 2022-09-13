package test_helpers

import (
	"encoding/json"
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func AssertStatusCode(t testing.TB, got *httptest.ResponseRecorder, want int) {
	t.Helper()
	if !Assert(t, got.Result().StatusCode, want, "response status code") {
		t.Fatalf("response: %v", got.Body)
	}
}

func Assert[T any](t testing.TB, got, want T, description string) bool {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s is not right:\ngot '%+v',\nwant '%+v'", description, got, want)
		return false
	}
	return true
}
func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected no error but got %v", got)
	}
}
func AssertError(t testing.TB, got error, want error) {
	t.Helper()
	if got != want {
		t.Errorf("expected error %v, but got %v", want, got)
	}
}
func AssertSomeError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Error("expected an error, but got nil")
	}
}

func AssertFatal[T comparable](t testing.TB, got, want T, description string) {
	t.Helper()
	if !Assert(t, got, want, description) {
		t.Fatal()
	}
}

func AssertClientError(t testing.TB, response *httptest.ResponseRecorder, err client_errors.ClientError) {
	t.Helper()
	var got client_errors.ClientError
	json.NewDecoder(response.Body).Decode(&got)

	AssertJSON(t, response)
	Assert(t, got.DisplayMessage, err.DisplayMessage, "detail code")
	Assert(t, response.Code, err.HTTPCode, "status code")
}

func AssertJSON(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	Assert(t, response.Result().Header.Get("contentType"), "application/json", "response content type")
}

func AssertJSONData[T any](t testing.TB, response *httptest.ResponseRecorder, wantData T) {
	t.Helper()
	AssertStatusCode(t, response, http.StatusOK)
	AssertJSON(t, response)
	var gotData T
	json.NewDecoder(response.Body).Decode(&gotData)
	Assert(t, gotData, wantData, "json encoded data")
}
