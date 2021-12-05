package gitsync

import (
	"fmt"
	"os"
	"testing"

	"github.com/blainemoser/gitsync/utils"
)

func TestMain(m *testing.M) {
	InitialiseTest()
	code := m.Run()
	TearDownTest()
	os.Exit(code)
}

func TestNewGit(t *testing.T) {
	if TestQueue == nil {
		t.Error("test queue instance could not be set")
	}
}

func TestStatus(t *testing.T) {
	errs := make([]error, 0)
	for _, process := range *TestQueue {
		result, err := process.Git().Status()
		if err != nil {
			errs = append(errs, err)
		}
		fmt.Println(result)
	}
	if len(errs) > 0 {
		t.Error(utils.ParseErrors(errs))
	}
}

func TestSync(t *testing.T) {
	result, err := SyncFile("two")
	fmt.Println(result)
	if err != nil {
		t.Error(err)
	}
	result, err = RemoveFileAndSync("two")
	fmt.Println(result)
	if err != nil {
		t.Error(err)
	}
}
