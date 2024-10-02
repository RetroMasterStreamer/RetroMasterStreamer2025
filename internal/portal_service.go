// internal/usecase.go
package internal

import (
	"PortalCRG/internal/repository/entity"
)

// PortalRetroGamerService representa los casos de uso relacionados con los usuarios.
type PortalRetroGamerService interface {
	Greet() string
	UpdateUserAvatar() string
	UpdateVideosTeams(search string) bool
	AuthenticateUser(alias, password string) (*entity.User, error)
	SetStatusLogin(alias, sessionToken, hash string, online bool) (bool, error)
	GetStatusLogin(sessionToken, hash string) (*entity.UserOnline, error)
	GetUserByAlias(alias string) (*entity.User, error)
	GetUserByTextRefer(text string) (*entity.User, error)
	ChangePassword(alias, password string) (*entity.User, error)
	SaveUser(user entity.User) (*entity.User, error)
	CreateUser(user *entity.User) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByRefer(refer string) (*entity.User, error)

	CreateTips(tips *entity.PostNew) (error, string)
	GetTipByID(id string) *entity.PostNew
	GetTipByURL(url string) *entity.PostNew

	GetTipsWithPagination(skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error)

	GetTipsByAliasWithPagination(alias string, skip, limit int64) ([]*entity.PostNew, error)

	GetTipsWithSearch(search string, skip, limit int64, typeOfTips []string) ([]*entity.PostNew, error)
	GetAllTips() ([]*entity.PostNew, error)
	DeleteTip(alias, id string) error
}
