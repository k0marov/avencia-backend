package client_errors

type ClientError struct {
	DetailCode string `json:"detail_code"`
	HTTPCode   int    `json:"-"`
}

var InvalidAuthToken = ClientError{
	DetailCode: "invalid-auth-token",
	HTTPCode:   401,
}

var InvalidJWT = ClientError{
	DetailCode: "invalid-jwt",
	HTTPCode:   400,
}
