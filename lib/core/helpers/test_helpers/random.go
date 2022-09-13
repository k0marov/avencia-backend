package test_helpers

import (
	"errors"
	"time"

	"github.com/AvenciaLab/avencia-api-contract/api"
	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core"
	authEntities "github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/entities"
	limitsEntities "github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/entities"
	limitsValues "github.com/AvenciaLab/avencia-backend/lib/features/limits/domain/values"
	transValues "github.com/AvenciaLab/avencia-backend/lib/features/transactions/domain/values"
	transferValues "github.com/AvenciaLab/avencia-backend/lib/features/transfers/domain/values"
	userEntities "github.com/AvenciaLab/avencia-backend/lib/features/users/domain/entities"
	walletEntities "github.com/AvenciaLab/avencia-backend/lib/features/wallets/domain/entities"
)


func RandomMetaTrans() transValues.MetaTrans {
	return transValues.MetaTrans{
		Type: RandomTransactionType(),
		UserId:    RandomString(),
	}
}

func RandomGeneratedCode() transValues.GeneratedCode {
	return transValues.GeneratedCode{
		Code:      RandomString(),
		ExpiresAt: time.Now(),
	}
}

func RandomClientError() client_errors.ClientError {
	return client_errors.ClientError{
		DisplayMessage: RandomString(),
		HTTPCode:   400 + RandomInt(),
	}
}

func RandomTransfer() transferValues.Transfer {
	return transferValues.Transfer{
		FromId: RandomString(),
		ToId:   RandomString(),
		Money:  RandomPositiveMoney(),
	}
}

func RandomRawTransfer() transferValues.RawTransfer {
	return transferValues.RawTransfer{
		FromId:  RandomString(),
		ToEmail: RandomString(),
		Money:   RandomPositiveMoney(),
	}
}

func RandomAPIMoney() api.Money {
	return api.Money{
		Currency: RandomString(),
		Amount:   RandomMoneyAmount().Num(),
	}
}

func RandomUser() authEntities.User {
	return authEntities.User{Id: RandomString()}
}
func RandomUserInfo() userEntities.UserInfo {
	return userEntities.UserInfo{Id: RandomString(), Wallet: RandomWallet(), Limits: RandomLimits()}
}
func RandomTransactionData() transValues.Transaction {
	return transValues.Transaction{
		Source: RandomTransactionSource(),
		UserId: RandomString(),
		Money:  RandomPositiveMoney(),
	}
}

func RandomTransactionSource() transValues.TransSource {
	return transValues.TransSource{
		Type:   transValues.TransSourceType(RandomString()),
		Detail: RandomString(),
	}
}

func RandomSecret() []byte {
	return []byte(RandomString())
}

func RandomWallet() walletEntities.Wallet {
	return walletEntities.Wallet{}
}

func RandomCurrency() core.Currency {
	return core.Currency(RandomString())
}

func RandomPosMoneyAmount() core.MoneyAmount {
	return core.NewMoneyAmount(RandomFloat())
}
func RandomNegMoneyAmount() core.MoneyAmount {
	return core.NewMoneyAmount(-RandomFloat())
}

func RandomMoneyAmount() core.MoneyAmount {
	if RandomBool() {
		return RandomPosMoneyAmount()
	} else {
		return RandomNegMoneyAmount()
	}
}

func RandomPositiveMoney() core.Money {
	return core.Money{
		Currency: RandomCurrency(),
		Amount:   RandomPosMoneyAmount(),
	}
}
func RandomNegativeMoney() core.Money {
	return core.Money{
		Currency: RandomCurrency(),
		Amount:   RandomNegMoneyAmount(),
	}
}

func RandomMoney() core.Money {
	return core.Money{
		Currency: RandomCurrency(), 
		Amount: RandomMoneyAmount(),
	}
}

func RandomLimits() limitsEntities.Limits {
	return limitsEntities.Limits{RandomCurrency(): RandomLimit(), RandomCurrency(): RandomLimit()}
}

func RandomLimit() limitsValues.Limit {
	return limitsValues.Limit{
		Withdrawn: RandomPosMoneyAmount(),
		Max:       RandomPosMoneyAmount(),
	}
}


func RandomTransactionType() transValues.TransactionType {
	if RandomBool() {
		return transValues.Deposit
	} else {
		return transValues.Withdrawal
	}
}

func RandomError() error {
	return errors.New(RandomString())
}

func RandomBool() bool {
	return randGen.Float32() > 0.5
}

func RandomInt() int {
	return randGen.Intn(100)
}
func RandomFloat() float64 {
	return randGen.Float64() * 100
}

func RandomString() string {
	str := ""
	for i := 0; i < 2; i++ {
		str += words[randGen.Intn(len(words))] + "_"
	}
	return str
}

var words = []string{"the", "be", "to", "of", "and", "a", "in", "that", "have", "I", "it", "for", "not", "on", "with", "he", "as", "you", "do", "at", "this", "but", "his", "by", "from", "they", "we", "say", "her", "she", "or", "an", "will", "my", "one", "all", "would", "there", "their", "what", "so", "up", "out", "if", "about", "who", "get", "which", "go", "me", "when", "make", "can", "like", "time", "no", "just", "him", "know", "take", "people", "into", "year", "your", "good", "some", "could", "them", "see", "other", "than", "then", "now", "look", "only", "come", "its", "over", "think", "also", "back", "after", "use", "two", "how", "our", "work", "first", "well", "way", "even", "new", "want", "because", "any", "these", "give", "day", "most", "us"}
