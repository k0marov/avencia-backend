package integration_test

import (
	"fmt"
	"testing"

	. "github.com/k0marov/avencia-backend/lib/core/helpers/test_helpers"
	"github.com/k0marov/avencia-backend/lib/di"
)

// TODO: maybe stop testing the handlers layer in here, but just call services directly

const (
	AtmAuthSecret = "atm_test"
	JwtSecret     = "jwt_test"
)

func TestIntegration(t *testing.T) {
	users := []MockUser{
		{
			Token: RandomString(),
			Id:    "sam",
			Email: "sam@skomarov.com",
		},
		{
			RandomString(),
			"john",
			"test@example.com",
		},
		{
			RandomString(),
			"bill",
			"test2@example.com",
		},
	}
	extDeps, cancelDBTrans := prepareExternalDeps(t, users)
	defer cancelDBTrans()
	initApiDeps(di.InitializeBusiness(extDeps), AtmAuthSecret, JwtSecret)

	fmt.Println("assert that balance of user 1 is 0$")
	assertBalance(t, users[0], newMoney("USD", 0))
	fmt.Println("assert that balance of user 2 is 0$")
	assertBalance(t, users[1], newMoney("USD", 0))
	fmt.Println("deposit 100$ to user 1")
	deposit(t, users[0], newMoney("USD", 100))
	fmt.Println("assert that balance of user 1 is 100$")
	assertBalance(t, users[0], newMoney("USD", 100))
	fmt.Println("withdraw 49.5$ from user1")
	withdraw(t, users[0], newMoney("USD", 49.5))
	fmt.Println("assert that balance of user1 is 50.5$")
	assertBalance(t, users[0], newMoney("USD", 50.5))
	fmt.Println("transfer 10.5$ from user1 to user2")
	transfer(t, users[0], users[1], newMoney("USD", 10.5))
	fmt.Println("assert that balance of user1 is 40$")
	assertBalance(t, users[0], newMoney("USD", 40))
	fmt.Println("assert that balance of user2 is 10.5$")
	assertBalance(t, users[1], newMoney("USD", 10.5))
	fmt.Println("withdraw 4.2$ from user2")
	withdraw(t, users[1], newMoney("USD", 4.2))
	fmt.Println("assert that balance of user2 is 6.3$")
	assertBalance(t, users[1], newMoney("USD", 6.3))
	fmt.Println("deposit 5000 RUB to user2")
	deposit(t, users[1], newMoney("RUB", 5000))
	fmt.Println("assert that balance of user2 is 5000 RUB")
	assertBalance(t, users[1], newMoney("RUB", 5000))
}

