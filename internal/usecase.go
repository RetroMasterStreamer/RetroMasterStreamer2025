// internal/usecase.go
package internal

import "PortalCRG/internal/repository/entity"

// UserService representa los casos de uso relacionados con los usuarios.
type UserService interface {
	Greet() string
	AuthenticateUser(alias, password string) (*entity.User, error)
	SetStatusLogin(alias, sessionToken, hash string, online bool) (bool, error)
	GetStatusLogin(sessionToken, hash string) (*entity.UserOnline, error)
	GetUserByAlias(alias string) (*entity.User, error)
	ChangePassword(alias, password string) (*entity.User, error)
	SaveUser(user entity.User) (*entity.User, error)
}
