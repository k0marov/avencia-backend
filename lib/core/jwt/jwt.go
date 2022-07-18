package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/k0marov/avencia-backend/secrets"
	"time"
)

const jwtSecret = secrets.JwtSecret

type Issuer func(subject string, expDuration time.Duration) (string, error)
type Verifier func(token string) (map[string]any, error)

func IssuerImpl(subject string, expDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().UTC().Add(expDuration).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func VerifierImpl(tokenString string) (map[string]any, error) {
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
