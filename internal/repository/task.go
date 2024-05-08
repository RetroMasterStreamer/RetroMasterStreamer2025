package repository

import "PortalCRG/internal/repository/entity"

// AuthenticationRepository define los métodos para la autenticación de usuarios.
type AuthenticationRepository interface {
	AuthenticateUser(alias, password string) (*entity.User, error)
}
