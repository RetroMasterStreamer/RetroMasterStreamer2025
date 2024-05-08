// internal/usecase_impl.go
package internal

import (
	"PortalCRG/internal/repository"
	"PortalCRG/internal/repository/entity"
)

// UserServiceImpl es una implementación de UserService.
type UserServiceImpl struct {
	UserRepository repository.UserRepositoryMongo
}

// NewUserService crea una nueva instancia de UserServiceImpl.
func NewUserService(userRepository repository.UserRepositoryMongo) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

// Greet retorna un saludo simple.
func (s *UserServiceImpl) Greet() string {
	return "Hello, world!"
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (s *UserServiceImpl) AuthenticateUser(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.AuthenticateUser(alias, password)
	return usuarioOnline, err
}
