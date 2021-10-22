package repository_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ismtabo/sso-poc/auth/pkg/repository"
	"github.com/stretchr/testify/assert"
)

var (
	pagesPath string
	pageRepo  repository.PageRepository
)

func setupPageRepoTests(m *testing.M) {
	var err error
	pagesPath, err = ioutil.TempDir("/tmp", "sso-poc-tests-*")
	if err != nil {
		log.Fatal(err)
	}
	pageRepo = repository.NewPageRepository(pagesPath)
}

func cleanupPageRepoTests(m *testing.M) {
	os.Remove(pagesPath)
}
func TestPageRepoPage(t *testing.T) {
	expectedPageName := fmt.Sprintf("%s.html", t.Name())
	if _, err := os.Create(filepath.Join(pagesPath, expectedPageName)); err != nil {
		t.Error(err)
	}
	actualPage, err := pageRepo.Page(expectedPageName)
	assert.NoError(t, err)
	assert.NotNil(t, actualPage)
	assert.Equal(t, []byte{}, actualPage)
}

func TestPageRepoPageNotFound(t *testing.T) {
	_, err := pageRepo.Page("unknown.ext")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
