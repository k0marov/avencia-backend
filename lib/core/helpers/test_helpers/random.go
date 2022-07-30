package test_helpers

import (
	"errors"
	"github.com/k0marov/avencia-api-contract/api"
	"github.com/k0marov/avencia-api-contract/api/client_errors"
	"github.com/k0marov/avencia-backend/lib/core"
	transValues "github.com/k0marov/avencia-backend/lib/features/atm/domain/values"
	"github.com/k0marov/avencia-backend/lib/features/auth"
	limitsEntities "github.com/k0marov/avencia-backend/lib/features/limits/domain/entities"
	limitsValues "github.com/k0marov/avencia-backend/lib/features/limits/domain/values"
	transferValues "github.com/k0marov/avencia-backend/lib/features/transfers/domain/values"
	userEntities "github.com/k0marov/avencia-backend/lib/features/users/domain/entities"
	walletEntities "github.com/k0marov/avencia-backend/lib/features/wallets/domain/entities"
	"math/rand"
	"strconv"
)

func RandomClientError() client_errors.ClientError {
	return client_errors.ClientError{
		DetailCode: RandomString(),
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

func RandomTransferRequest() api.TransferRequest {
	return api.TransferRequest{
		RecipientIdentifier: RandomString(),
		Currency:            RandomString(),
		Amount:              RandomFloat(),
	}
}

func RandomUser() auth.User {
	return auth.User{Id: RandomId()}
}
func RandomUserInfo() userEntities.UserInfo {
	return userEntities.UserInfo{Id: RandomId(), Wallet: RandomWallet(), Limits: RandomLimits()}
}
func RandomTransactionData() transValues.Transaction {
	return transValues.Transaction{
		UserId: RandomString(),
		Money:  RandomPositiveMoney(),
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

func RandomBanknoteCheckRequest() api.BanknoteCheckRequest {
	return api.BanknoteCheckRequest{
		TransactionCode: RandomString(),
		Currency:        RandomString(),
		Amount:          RandomFloat(),
	}
}

func RandomFinalizeTransationRequest() api.FinalizeTransactionRequest {
	return api.FinalizeTransactionRequest{
		UserId:    RandomString(),
		ATMSecret: RandomString(),
		Currency:  RandomString(),
		Amount:    RandomFloat(),
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

func RandomBanknote() transValues.Banknote {
	return transValues.Banknote{Money: RandomPositiveMoney()}
}

func RandomTransactionType() transValues.TransactionType {
	if RandomBool() {
		return transValues.Deposit
	} else {
		return transValues.Withdrawal
	}
}

func RandomId() string {
	return strconv.Itoa(rand.Intn(100000))
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
