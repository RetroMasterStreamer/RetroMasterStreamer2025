// internal/usecase.go
package internal

import "PortalCRG/internal/repository/entity"

// UserService representa los casos de uso relacionados con los usuarios.
type UserService interface {
	Greet() string
	AuthenticateUser(alias, password string) (*entity.User, error)
}
