package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/delivery/http/middleware"
)

func TestAuthMiddleware(t *testing.T) {
	authHeader := RandomString()
	requestCtx := context.WithValue(context.Background(), RandomString(), RandomString())
	returnedCtx := context.WithValue(context.Background(), RandomString(), RandomString())

	r := httptest.NewRequest("", "/asdf", nil).WithContext(requestCtx)
	r.Header.Add("Authorization", authHeader)

	adder := func(gotCtx context.Context, gotHeader string) context.Context {
		if gotCtx == requestCtx && gotHeader == authHeader {
			return returnedCtx
		}
		panic("unexpected")
	}

	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		Assert(t, r.Context(), returnedCtx, "added ctx")
	}

	mw := middleware.NewAuthMiddleware(adder)
	mw(http.HandlerFunc(handler)).ServeHTTP(httptest.NewRecorder(), r)
	Assert(t, called, true, "next handler was called")
}
