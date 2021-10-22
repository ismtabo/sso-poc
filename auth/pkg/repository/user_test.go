package repository_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ismtabo/sso-poc/auth/pkg/model"
	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"github.com/ismtabo/sso-poc/auth/pkg/repository/dao"
	"github.com/sonyarouje/simdb"
	"github.com/stretchr/testify/assert"
)

var simdbPath string
var driver *simdb.Driver
var userRepo repository.UserRepository

func setupUserRepoTests(m *testing.M) {
	var err error
	simdbPath, err = ioutil.TempDir("/tmp", "sso-poc-test-*")
	if err != nil {
		log.Fatal(err)
	}
	driver, err = simdb.New(simdbPath)
	if err != nil {
		log.Fatal(err)
	}
	userRepo = repository.NewSimdbUserRepository(driver)
}

func cleanupUserRepoTests(m *testing.M) {
	defer os.Remove(simdbPath)
}

func TestUserRepositoryCreateUser(t *testing.T) {
	expectedUser := &model.User{UID: t.Name(), Password: "password"}
	actualUser, err := userRepo.CreateUser(expectedUser)
	assert.NoError(t, err)
	assert.NotNil(t, actualUser)
	assert.EqualValues(t, expectedUser, actualUser)
}

func TestUserRepositoryCreateUserDuplicatedUID(t *testing.T) {
	driver.Open(dao.User{}).Insert(&dao.User{UID: t.Name()})
	expectedUser := &model.User{UID: t.Name()}
	_, err := userRepo.CreateUser(expectedUser)
	assert.Error(t, err)
}

func TestUserRepositoryReadUser(t *testing.T) {
	driver.Open(dao.User{}).Insert(&dao.User{UID: t.Name()})
	expectedUser := &model.User{UID: t.Name()}
	actualUser, err := userRepo.ReadUser(expectedUser.UID)
	assert.NoError(t, err)
	assert.NotNil(t, actualUser)
	assert.EqualValues(t, expectedUser, actualUser)
}

func TestUserRepositoryReadUserNotFound(t *testing.T) {
	_, err := userRepo.ReadUser("unknown")
	assert.Error(t, err)
}

func TestUserRepositoryDeleteUser(t *testing.T) {
	driver.Open(dao.User{}).Insert(&dao.User{UID: t.Name()})
	expectedUser := &model.User{UID: t.Name()}
	err := userRepo.DeleteUser(expectedUser.UID)
	assert.NoError(t, err)
}

func TestUserRepositoryDeleteUserNotFound(t *testing.T) {
	err := userRepo.DeleteUser("unknown")
	assert.Error(t, err)
}
