package entities

import "github.com/k0marov/avencia-backend/api"

type UserInfo struct {
	Id string
}

func (u UserInfo) ToResponse() api.UserInfoResponse {
	return api.UserInfoResponse{Id: u.Id}
}
