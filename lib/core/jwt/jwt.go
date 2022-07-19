package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Issuer func(claims map[string]any, expireAt time.Time) (string, error)
type Verifier func(token string) (map[string]any, error)

func NewIssuer(jwtSecret []byte) Issuer {
	return func(claims map[string]any, expireAt time.Time) (string, error) {
		claims["exp"] = expireAt.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

		return token.SignedString(jwtSecret)
	}
}

func NewVerifier(jwtSecret []byte) Verifier {
	return func(tokenString string) (map[string]any, error) {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})
		if err != nil {
			return make(map[string]any), fmt.Errorf("while parsing a token: %w", err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		} else {
			return make(map[string]any), fmt.Errorf("token is invalid")
		}
	}
}
