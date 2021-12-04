package sync

import (
	"fmt"
	"io/fs"
	"os"

	jsonextract "github.com/blainemoser/JsonExtract"
)

var TestDir = "../../../testGitSync"

var TestGit *Git

// InitialiseTest initializes the testing environment
func InitialiseTest() {
	dir, err := checkTestDir()
	if err != nil {
		panic(err)
	}
	TestGit, err = NewGit().SetRepo(dir)
	if err != nil {
		panic(err)
	}
}

func TearDownTest() {
	if TestGit != nil {
		dir := TestGit.GetRepo()
		if dir == TestDir {
			err := os.RemoveAll(TestDir)
			if err != nil {
				panic(err)
			}
		}
	}
}

func checkTestDir() (string, error) {
	setup, err := os.OpenFile("../test.json", os.O_RDONLY, os.ModeDevice)
	if err != nil {
		if os.IsNotExist(err) {
			return makeTestDir()
		}
		return "", err
	}
	return setTestDir(setup)
}

func setTestDir(setup *os.File) (string, error) {
	stat, err := setup.Stat()
	if err != nil {
		return "", err
	}
	b := make([]byte, stat.Size())
	_, err = setup.Read(b)
	if err != nil {
		return "", err
	}
	return setTestRepo(b)
}

func setTestRepo(b []byte) (string, error) {
	// Decode the contents
	json := jsonextract.JSONExtract{RawJSON: string(b)}
	path, err := json.Extract("testRepo")
	if err != nil {
		return "", err
	}
	if result, ok := path.(string); ok {
		return result, nil
	}
	return "", fmt.Errorf("path to repo not found in test.json")
}

func makeTestDir() (string, error) {
	err := os.Mkdir(TestDir, fs.FileMode(0777))
	if os.IsExist(err) {
		fmt.Printf("warning: %s\n", err.Error())
		return "", nil
	}
	return TestDir, err
}

// SyncFile syncs a file to the repo
func SyncFile(name string) (string, error) {
	err := makeFile(name)
	if err != nil {
		return "", err
	}
	return TestGit.Sync()
}

// RemoveFileAndSync removes the file then syncs the git repo
func RemoveFileAndSync(name string) (string, error) {
	err := os.Remove(TestGit.GetRepo() + "/" + name + ".txt")
	if err != nil {
		return "", err
	}
	return TestGit.Sync()
}

func makeFile(name string) error {
	file, err := os.Create(TestGit.GetRepo() + "/" + name + ".txt")
	if err != nil {
		return err
	}
	content := []byte(TestFiles[name])
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}
