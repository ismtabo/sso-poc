package service

import (
	"github.com/ismtabo/sso-poc/auth/pkg/model"
	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(uid, password string) error
	Register(uid, password string) error
}

type authService struct {
	users repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{users: userRepo}
}

func (as *authService) Authenticate(uid, password string) error {
	user, err := as.users.ReadUser(uid)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
}

func (as *authService) Register(uid, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = as.users.CreateUser(&model.User{UID: uid, Password: string(hashedPwd)})
	if err != nil {
		return err
	}
	return nil
}
