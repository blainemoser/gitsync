package sync

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitialiseTest()
	code := m.Run()
	TearDownTest()
	os.Exit(code)
}

func TestNewGit(t *testing.T) {
	if TestGit == nil {
		t.Error("test git instance could not be set")
	}
	fmt.Println(TestGit.GetRepo())
}
