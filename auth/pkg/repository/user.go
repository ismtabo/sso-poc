package repository

import (
	"errors"
	"fmt"

	"github.com/ismtabo/sso-poc/auth/pkg/model"
	"github.com/ismtabo/sso-poc/auth/pkg/repository/dao"
	"github.com/jinzhu/copier"
	"github.com/sonyarouje/simdb"
)

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	ReadUser(uid string) (*model.User, error)
	DeleteUser(uid string) error
}

type simdbUserRepository struct {
	driver *simdb.Driver
}

func NewSimdbUserRepository(driver *simdb.Driver) UserRepository {
	return &simdbUserRepository{driver: driver}
}

func (ur *simdbUserRepository) CreateUser(user *model.User) (*model.User, error) {
	if otherUser, err := ur.ReadUser(user.UID); !errors.Is(err, simdb.ErrRecordNotFound) || otherUser != nil {
		if err != nil {
			return nil, err
		} else if otherUser != nil {
			return nil, fmt.Errorf("already exists user with id '%s'", user.UID)
		}
	}
	userDao := &dao.User{}
	copier.Copy(userDao, user)
	if err := ur.driver.Open(dao.User{}).Insert(userDao); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *simdbUserRepository) ReadUser(uid string) (*model.User, error) {
	userDao := &dao.User{}
	if err := ur.driver.Open(dao.User{}).Where("username", "=", uid).First().AsEntity(userDao); err != nil {
		return nil, err
	}
	user := &model.User{}
	copier.Copy(user, userDao)
	return user, nil
}

func (ur *simdbUserRepository) DeleteUser(uid string) error {
	return ur.driver.Open(dao.User{}).Delete(&dao.User{UID: uid})
}
