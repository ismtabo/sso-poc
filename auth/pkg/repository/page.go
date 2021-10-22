package repository

import (
	"io/ioutil"
	"path/filepath"
)

type PageRepository interface {
	Page(path string) ([]byte, error)
}

type pageRepository struct {
	pagesPath string
}

func NewPageRepository(pagesPath string) PageRepository {
	return &pageRepository{pagesPath: pagesPath}
}

func (s *pageRepository) Page(page string) ([]byte, error) {
	path := filepath.Join(s.pagesPath, page)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
