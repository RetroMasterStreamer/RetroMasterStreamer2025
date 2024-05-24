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

type NewUserRequest struct {
	NewUser entity.NewUser `json:"new"`
	RefUser entity.User    `json:"ref"`
	Code    string         `json:"code"`
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

type ResponseMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type TipsShared struct {
	Title    string
	ID       string
	URL      string
	Content  string
	AvatarYT string
}
