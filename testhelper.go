package main

import (
	"os"

	"github.com/blainemoser/gitsync/configs"
	"github.com/blainemoser/gitsync/queue"
	"github.com/blainemoser/gitsync/utils"
)

var TestDir = "../../../testGitSync"

var TestQueue *queue.Queue

// InitialiseTest initializes the testing environment
func initialiseTest() {
	pwd()
	configs, err := configs.NewConfigs().SetDirectories("test.json", TestDir)
	if err != nil {
		panic(err)
	}
	TestQueue, err = queue.NewQueue().Walk(configs)
	if err != nil {
		panic(err)
	}
}

func tearDownTest() {
	if TestQueue == nil {
		return
	}
	for _, process := range *TestQueue {
		if process.Git().GetRepo() != TestDir {
			continue
		}
		err := os.RemoveAll(TestDir)
		if err != nil {
			panic(err)
		}
		return
	}
}

func syncFile(errChan chan error, name string) {
	errs := make([]error, 0)
	var file *os.File
	var err error
	var content []byte
	for _, process := range *TestQueue {
		file, err = os.Create(process.Git().GetRepo() + "/" + name + ".txt")
		if err != nil {
			errs = append(errs, err)
			continue
		}
		content = []byte(TestFiles[name])
		_, err = file.Write(content)
		if err != nil {
			errs = append(errs, err)
		}
	}
	errChan <- utils.ParseErrors(errs)
}

func removeFileAndSync(errChan chan error, name string) {
	errs := make([]error, 0)
	var err error
	for _, process := range *TestQueue {
		err = os.Remove(process.Git().GetRepo() + "/" + name + ".txt")
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		errChan <- utils.ParseErrors(errs)
		return
	}
	errChan <- nil
}

func updateFileContent(errChan chan error, name string) {
	errs := make([]error, 0)
	var file *os.File
	var err error
	var content []byte
	for _, process := range *TestQueue {
		file, err = os.OpenFile(process.Git().GetRepo()+"/"+name+".txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		content = []byte("\nUPDATE\nNew Write to File\n")
		_, err = file.Write(content)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		errChan <- utils.ParseErrors(errs)
		return
	}
	errChan <- nil
}
