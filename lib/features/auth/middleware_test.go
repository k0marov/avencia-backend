package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth"
)

func TestMiddleware(t *testing.T) {

	assertResult := func(verify auth.Verifier, r *http.Request, wantUser auth.User) {
    called := false 
    handler := func(w http.ResponseWriter, r *http.Request) {
    	called = true
      user, err := auth.UserFromCtx(r.Context())
      noUser := auth.User{}
      if wantUser != noUser {
      	Assert(t, user, wantUser, "")
      	AssertNoError(t, err)
      } else {
      	AssertError(t, err, auth.ErrNoUserInContext)
      }
    } 
    auth.NewFirebaseAuthMiddleware(verify)(http.HandlerFunc(handler)).ServeHTTP(httptest.NewRecorder(), r) 
    Assert(t, called, true, "next handler was called")
	}

  t.Run("error case - token is not provided", func(t *testing.T) {
    request := httptest.NewRequest("", "/asdf", nil)
    assertResult(nil, request, auth.User{})
  }) 
  token := RandomString() 
	
	genGoodReq := func() *http.Request {
		request := httptest.NewRequest("", "/asdf", nil) 
		request.Header.Add("Authorization", "Bearer " + token)
		return request
	}

  t.Run("error case - token is invalid", func(t *testing.T) {
		verifier := func(string) (string, bool) {
			return "", false 
		}
		assertResult(verifier, genGoodReq(), auth.User{})
  })
	
  t.Run("happy case", func(t *testing.T) {
  	user := RandomUser() 
  	verifier := func(gotToken string) (string, bool) {
  		if gotToken == token {
  			return user.Id, true
  		} 
  		panic("unexpected")
  	}
  	assertResult(verifier, genGoodReq(), user) 
  })
}
