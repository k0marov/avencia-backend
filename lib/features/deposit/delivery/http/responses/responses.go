package responses

type UserInfoResponse struct {
	Id string `json:"id"`
}

type CodeResponse struct {
	TransactionCode string `json:"transaction_code"`
	ExpiresAt       int64  `json:"expires_at"`
}

type BanknoteCheckResponse struct {
	Accept bool `json:"accept"`
}
