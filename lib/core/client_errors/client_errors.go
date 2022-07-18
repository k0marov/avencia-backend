package client_errors

import "fmt"

type ClientError struct {
	DetailCode string `json:"detail_code"`
	HTTPCode   int    `json:"-"`
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("An error which will be displayed to the client: %v %v", ce.HTTPCode, ce.DetailCode)
}

var InvalidAuthToken = ClientError{
	DetailCode: "invalid-auth-token",
	HTTPCode:   401,
}

var InvalidJWT = ClientError{
	DetailCode: "invalid-jwt",
	HTTPCode:   400,
}
