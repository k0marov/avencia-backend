package api

type UserInfoResponse struct {
	Id     string             `json:"id"`
	Wallet map[string]float64 `json:"wallet"`
}
type VerifiedCodeResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}

type CodeResponse struct {
	TransactionCode string `json:"transaction_code"`
	ExpiresAt       int64  `json:"expires_at"`
}
