package server

import "PortalCRG/internal/repository/entity"

type Credentials struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Password           string `json:"password"`
	NewPassword        string `json:"password_new"`
	ConfirmNewPassword string `json:"password_confirm_new"`
}

type ResponseLogin struct {
	User entity.User
	Hash string
}

type ResponseOnline struct {
	Status string
	Code   int
	User   entity.UserOnline
}
