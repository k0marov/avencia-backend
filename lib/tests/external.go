package integration_test

import (
	"testing"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/AvenciaLab/avencia-backend/lib/core/core_err"
	"github.com/AvenciaLab/avencia-backend/lib/core/db/foundationdb"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/general_helpers"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
)

type MockUser struct {
	// User should be provided with at least the Id and the Email
	User authEntities.DetailedUser 
	Token string
}

type MockAuth struct {
	Users []MockUser
}

func (a MockAuth) Verify(token string) (userId string) {
	for _, u := range a.Users {
		if u.Token == token {
			return u.User.Id
		}
	}
	return ""
}

func (a MockAuth) UserByEmail(email string) (authEntities.User, error) {
	for _, u := range a.Users {
		if u.User.Email == email {
			return authEntities.User{Id: u.User.Id}, nil
		}
	}
	return authEntities.User{}, core_err.ErrNotFound
}

func (a MockAuth) Get(id string) (authEntities.DetailedUser, error) {
	for _, u := range a.Users {
		if u.User.Id == id {
			return u.User, nil
		}
	}
	return authEntities.DetailedUser{}, core_err.ErrNotFound
}


func prepareExternalDeps(t *testing.T, users []MockUser) (d di.ExternalDeps, cancelTrans func()) {
	fdb.MustAPIVersion(710) 
	db := fdb.MustOpenDefault()
	// executing all the tests inside a transaction which will then be canceled
	trans, err := db.CreateTransaction()
	AssertNoError(t, err)
	trans.ClearRange(general_helpers.ConvTuple([]string{"", "\xFF"}))
	runner := foundationdb.NewTransactionRunner(trans)
	

	return di.ExternalDeps{
		AtmSecret: []byte("atm_test"),
		JwtSecret: []byte("jwt_test"),
		Auth:      MockAuth{Users: users},
		TRunner:   runner,
	}, trans.Cancel
}
