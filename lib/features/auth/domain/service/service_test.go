package service_test

import (
	"context"
	"testing"

	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
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
