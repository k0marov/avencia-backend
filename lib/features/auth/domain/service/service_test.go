package service_test

import (
	"testing"

	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/features/auth/domain/service"
	"golang.org/x/net/context"
)

func TestUserInfoAdder(t *testing.T) {
	token := RandomString()
	header := token
	t.Run("error case - the provided token is invalid", func(t *testing.T) {
		ctx := context.Background()
		verify := func(gotToken string) string {
			if gotToken == token {
				return ""
			}
			panic("unexpected")
		}
		gotCtx := service.NewUserInfoAdder(verify)(ctx, header)
		Assert(t, gotCtx, ctx, "returned context")
	})

	t.Run("happy case", func(t *testing.T) {
		user := RandomUser()
		ctx := context.Background()
		verify := func(string) string {
			return user.Id
		}
		gotCtx := service.NewUserInfoAdder(verify)(ctx, header)
		gotUser, err := service.UserFromCtx(gotCtx)
		AssertNoError(t, err)
		Assert(t, gotUser, user, "user that was put in ctx")
	})
}
