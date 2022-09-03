package tests

import (
	"testing"

	"github.com/k0marov/avencia-backend/lib/di"
)

func prepareExternalDeps() di.ExternalDeps {
  return di.ExternalDeps{
  	AtmSecret: []byte("atm_test"),
  	JwtSecret: []byte("jwt_test"),
  	Auth:      nil,
  	TRunner:   nil,
  }
} 

func TestIntegration(t *testing.T) {

}
