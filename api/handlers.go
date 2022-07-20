package api

import "net/http"

type ATMTransaction struct {
	GenCode, VerifyCode, CheckBanknote, FinalizeTransaction http.HandlerFunc
}
