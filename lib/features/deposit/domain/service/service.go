package service

import (
	"github.com/k0marov/avencia-backend/lib/core/client_errors"
	"github.com/k0marov/avencia-backend/lib/core/jwt"
	"github.com/k0marov/avencia-backend/lib/features/deposit/domain/entities"
)

type CodeVerifier = func(string) (entities.UserInfo, error)

func NewCodeVerifier(verifyJWT jwt.Verifier) CodeVerifier {
	return func(code string) (entities.UserInfo, error) {
		data, err := verifyJWT(code)
		if err != nil {
			return entities.UserInfo{}, client_errors.InvalidJWT
		}
		userId, _ := data["sub"].(string)
		return entities.UserInfo{
			Id: userId,
		}, nil
	}
}
