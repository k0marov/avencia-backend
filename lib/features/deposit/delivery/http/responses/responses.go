package responses

type UserInfoResponse struct {
	Id string `json:"id"`
}

type CodeResponse struct {
	TransactionCode string `json:"transaction_code"`
}
