package tests

import (
	"testing"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/k0marov/avencia-backend/lib/core/core_err"
	"github.com/k0marov/avencia-backend/lib/core/db/foundationdb"
	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/di"
	authEntities "github.com/k0marov/avencia-backend/lib/features/auth/domain/entities"
)

type MockUser struct {
	Token, Id, Email string 
}

type MockAuth struct {
  Users []MockUser 
}

func (a MockAuth) Verify(token string) (userId string) {
	for _, u := range a.Users {
		if u.Token == token {
			return u.Id
		}
	}
	return "" 
} 

func (a MockAuth) UserByEmail(email string) (authEntities.User, error) {
  for _, u := range a.Users {
  	if u.Email == email {
  		return authEntities.User{Id: u.Id}, nil
  	}
  }
  return authEntities.User{}, core_err.ErrNotFound
}

func prepareExternalDeps(t *testing.T, users []MockUser) di.ExternalDeps {
	db := fdb.MustOpenDefault()
	trans, err := db.CreateTransaction()
	AssertNoError(t, err)
	trans.SetReadVersion(1) 
  
  runner := foundationdb.NewTransactionRunner(trans)

  return di.ExternalDeps{
  	AtmSecret: []byte("atm_test"),
  	JwtSecret: []byte("jwt_test"),
  	Auth:      MockAuth{Users: users},
  	TRunner:   runner,
  }
} 


// TODO: maybe move middleware usage from the api repository here

func TestIntegration(t *testing.T) {
	users := []MockUser{
		{
			Token: RandomString(),
			Id:    "sam",
			Email: "sam@skomarov.com",
		}, 
		{
			Token: RandomString(),
			Id:    "john",
			Email: "test@example.com",
		}, 
		{
			Token:RandomString(),
			Id:    "bill",
			Email: "test2@example.com",
		}, 
	}
	apiDeps := di.InitializeBusiness(prepareExternalDeps(t, users))

	


}
