package repository_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setupUserRepoTests(m)
	setupPageRepoTests(m)
	code := m.Run()
	cleanupUserRepoTests(m)
	cleanupPageRepoTests(m)
	os.Exit(code)
}
