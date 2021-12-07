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
	TestQueue.StandbyAll()
	err := testFileCreate()
	if err != nil {
		t.Error(err)
	}
	err = testRemoveFile()
	if err != nil {
		t.Error(err)
	}
}

func testFileCreate() error {
	fmt.Println("creating test file")
	errChan := make(chan error, 1)
	fmt.Println("starting to create a file")
	SyncFile(errChan, "two")
	err := <-errChan
	close(errChan)
	if err != nil {
		return err
	}
	TestQueue.WaitForAll()
	return nil
}

func testRemoveFile() error {
	fmt.Println("removing a file")
	errChan := make(chan error, 1)
	RemoveFileAndSync(errChan, "two")
	err := <-errChan
	close(errChan)
	if err != nil {
		return err
	}
	TestQueue.WaitForAll()
	return nil
}
