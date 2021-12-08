package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/blainemoser/gitsync/utils"
)

func TestMain(m *testing.M) {
	initialiseTest()
	code := m.Run()
	tearDownTest()
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
	err := testFileCreate("two")
	if err != nil {
		t.Error(err)
	}
	err = testFileUpdate()
	if err != nil {
		t.Error(err)
	}
	err = testRemoveFile("two")
	if err != nil {
		t.Error(err)
	}
	err = testRemoveFile("one")
	if err != nil {
		t.Error(err)
	}
}

func testFileCreate(name string) error {
	fmt.Println("creating test file")
	errChan := make(chan error, 1)
	syncFile(errChan, name)
	err := <-errChan
	close(errChan)
	if err != nil {
		return err
	}
	TestQueue.WaitForAll()
	return nil
}

func testFileUpdate() error {
	fmt.Println("creating test file")
	err := testFileCreate("one")
	if err != nil {
		return err
	}
	errChan := make(chan error, 1)
	fmt.Println("updating test file")
	updateFileContent(errChan, "one")
	err = <-errChan
	close(errChan)
	if err != nil {
		return err
	}
	TestQueue.WaitForAll()
	return nil
}

func testRemoveFile(name string) error {
	fmt.Println("removing a file")
	errChan := make(chan error, 1)
	removeFileAndSync(errChan, name)
	err := <-errChan
	close(errChan)
	if err != nil {
		return err
	}
	TestQueue.WaitForAll()
	return nil
}
