package configs

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/blainemoser/gitsync/utils"
)

var TestDirName string

func TestMain(m *testing.M) {
	InitialiseTest()
	code := m.Run()
	TearDownTest()
	os.Exit(code)
}

func InitialiseTest() {
	TestDirName = fmt.Sprintf("../config_test_%s", time.Now().Format(time.RFC3339))
	err := makeFolder()
	if err != nil {
		panic(err)
	}
	err = makeTestSetup()
	if err != nil {
		panic(err)
	}
}

func TearDownTest() {
	err := removeTestSetup()
	if err != nil {
		panic(err)
	}
}

func TestDirectories(t *testing.T) {
	conf, err := NewConfigs().SetDirectories("../configs_test.json", "")
	if err != nil {
		t.Error(err)
	}
	directories := conf.GetDirectories()
	if len(directories) < 1 {
		t.Errorf("expected one directory to be found, found %d", len(directories))
	}
	if directories[0] != TestDirName {
		t.Errorf("expected directory to be called %s, got %s", TestDirName, directories[0])
	}
}

func makeFolder() error {
	err := os.Mkdir(TestDirName, os.FileMode(777))
	if err != nil {
		return err
	}
	return nil
}

func makeTestSetup() error {
	setup, err := os.Create("../configs_test.json")
	if err != nil {
		return err
	}
	content := []byte(`{"directories": ["` + TestDirName + `"]}`)
	_, err = setup.Write(content)
	return err
}

func removeTestSetup() error {
	var err error
	errs := make([]error, 0)
	err = os.RemoveAll(TestDirName)
	if err != nil {
		errs = append(errs, err)
	}
	err = os.Remove("../configs_test.json")
	if err != nil {
		errs = append(errs, err)
	}
	return utils.ParseErrors(errs)
}
