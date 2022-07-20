package api

type UserInfoResponse struct {
	Id string `json:"id"`
}
type VerifiedCodeResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}

type CodeResponse struct {
	TransactionCode string `json:"transaction_code"`
	ExpiresAt       int64  `json:"expires_at"`
}
