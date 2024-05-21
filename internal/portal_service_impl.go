// internal/usecase_impl.go
package internal

import (
	"PortalCRG/internal/repository"
	"PortalCRG/internal/repository/entity"
)

// PortalRetroGamerImpl es una implementación de UserService.
type PortalRetroGamerImpl struct {
	UserRepository   repository.UserRepositoryMongo
	PortalRepository repository.PortalRepositoryMongo
}

// NewUserService crea una nueva instancia de UserServiceImpl.
func NewUserService(userRepository repository.UserRepositoryMongo, portalRepository repository.PortalRepositoryMongo) *PortalRetroGamerImpl {
	return &PortalRetroGamerImpl{
		UserRepository:   userRepository,
		PortalRepository: portalRepository,
	}
}

// Greet retorna un saludo simple.
func (s *PortalRetroGamerImpl) Greet() string {
	return "Hello, world!"
}

// AuthenticateUser autentica a un usuario utilizando su alias y contraseña.
func (s *PortalRetroGamerImpl) AuthenticateUser(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.AuthenticateUser(alias, password)
	return usuarioOnline, err
}

func (s *PortalRetroGamerImpl) SetStatusLogin(alias, sessionToken, hash string, online bool) (bool, error) {
	usuarioOnline, err := s.UserRepository.SetUserOnline(alias, sessionToken, hash, online)
	if online {
		return usuarioOnline.Online, err
	} else {
		return false, err
	}
}

func (s *PortalRetroGamerImpl) GetStatusLogin(sessionToken, hash string) (*entity.UserOnline, error) {
	usuarioOnline, err := s.UserRepository.GetUserOnline(sessionToken, hash)
	return usuarioOnline, err
}
func (s *PortalRetroGamerImpl) GetUserByAlias(alias string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	return usuarioOnline, err
}
func (s *PortalRetroGamerImpl) GetUserByTextRefer(text string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByTextRefer(text)
	return usuarioOnline, err
}

func (s *PortalRetroGamerImpl) ChangePassword(alias, password string) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Password = password
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}

func (s *PortalRetroGamerImpl) SaveUser(user entity.User) (*entity.User, error) {
	usuarioOnline, err := s.UserRepository.GetUserByAlias(user.Alias)
	if err != nil {
		return nil, err
	} else {
		usuarioOnline.Name = user.Name
		usuarioOnline.RRSS = user.RRSS
		usuarioOnline.AvatarYT = user.AvatarYT
		usuarioOnline.ReferenceText = user.ReferenceText
		usuarioOnline.AboutMe = user.AboutMe
		s.UserRepository.SaveUser(usuarioOnline)
		return usuarioOnline, nil
	}
}

func (s *PortalRetroGamerImpl) CreateUser(user *entity.User) error {
	error := s.UserRepository.SaveUser(user)
	return error
}

func (s *PortalRetroGamerImpl) GetAllUsers() ([]*entity.User, error) {

	users, err := s.PortalRepository.GetAllUsers()

	return users, err
}

func (s *PortalRetroGamerImpl) GetUserByRefer(refer string) (*entity.User, error) {

	user, err := s.PortalRepository.GetUserByAlias(refer)

	return user, err
}

func (s *PortalRetroGamerImpl) GetAllTips() ([]*entity.PostNew, error) {

	users, err := s.PortalRepository.GetAllTips()

	return users, err
}

func (s *PortalRetroGamerImpl) CreateTips(tip *entity.PostNew) error {
	error := s.UserRepository.SaveTips(tip)
	return error
}

func (s *PortalRetroGamerImpl) GetTipByID(id string) *entity.PostNew {

	tips, _ := s.PortalRepository.GetTipByID(id)
	return tips
}

func (s *PortalRetroGamerImpl) GetTipsWithPagination(skip, limit int64) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsWithPagination(skip, limit)
}

func (s *PortalRetroGamerImpl) GetTipsWithSearch(search string, skip, limit int64) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsWithSearch(search, skip, limit)
}

func (s *PortalRetroGamerImpl) DeleteTip(id, alias string) error {
	return s.PortalRepository.DeleteTipByIDandAuthor(id, alias)
}

func (s *PortalRetroGamerImpl) GetTipsByAliasWithPagination(alias string, skip, limit int64) ([]*entity.PostNew, error) {
	return s.PortalRepository.GetTipsByAliasWithPagination(alias, skip, limit)
}
