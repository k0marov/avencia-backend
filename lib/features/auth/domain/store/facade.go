package facade

// AuthFacade Verify returns "" if the provided token is invalid
type AuthFacade interface {
	Verify(token string) (userId string)
}

