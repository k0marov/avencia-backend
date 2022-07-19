package api

import "net/http"

type CashDeposit struct {
	GenCode, VerifyCode, CheckBanknote, FinalizeTransaction http.HandlerFunc
}
