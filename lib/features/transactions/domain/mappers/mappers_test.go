package mappers_test

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/AvenciaLab/avencia-backend/lib/setup/config/configurable"
	. "github.com/AvenciaLab/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/core/jwt"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/mappers"
	"github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
)

func TestPropertiesOfCodeMapping(t *testing.T) {
	jwtSecret := RandomSecret()
	generate := mappers.NewCodeGenerator(jwt.NewIssuer(jwtSecret))
	parse := mappers.NewCodeParser(jwt.NewVerifier(jwtSecret))
	assertion := func(trans values.MetaTrans) bool {
		code, err := generate(trans)
		if err != nil {
			return false 
		}

		wantExpireAt := time.Now().UTC().Add(configurable.TransactionExpDuration)
		Assert(t, TimeAlmostEqual(code.ExpiresAt, wantExpireAt), true, "the expiration time is Now + ExpDuration")

		parsedTrans, err := parse(code.Code) 
		return parsedTrans == trans && err == nil
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Error("failed checks", err)
	}
}

func TestPropertiesOfTransactionIdEncoding(t *testing.T) {
	generate := mappers.NewTransIdGenerator()
	parse := mappers.NewTransIdParser()
	assertion := func(code string) bool {
		return parse(generate(code)) == code
	}

	if err := quick.Check(assertion, nil); err != nil {
		t.Error("failed checks", err)
	}
}

