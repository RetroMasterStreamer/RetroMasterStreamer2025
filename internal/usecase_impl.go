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

func (s *UserServiceImpl) SetStatusLogin(alias, sessionToken, hash string, online bool) (bool, error) {
	usuarioOnline, err := s.UserRepository.SetUserOnline(alias, sessionToken, hash, online)
	if online {
		return usuarioOnline.Online, err
	} else {
		return false, err
	}
}

func (s *UserServiceImpl) GetStatusLogin(sessionToken, hash string) (*entity.UserOnline, error) {
	usuarioOnline, err := s.UserRepository.GetUserOnline(sessionToken, hash)
	return usuarioOnline, err
}
func (s *UserServiceImpl) GetUserByAlias(alias string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	return usuarioOnline, err
}

func (s *UserServiceImpl) ChangePassword(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Password = password
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}

func (s *UserServiceImpl) SaveUser(user entity.User) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(user.Alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Name = user.Name
		usuarioOnline.RRSS = user.RRSS
		usuarioOnline.ReferenceText = user.ReferenceText
		usuarioOnline.AboutMe = user.AboutMe
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}
