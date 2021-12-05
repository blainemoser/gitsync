package configs

import (
	"fmt"
	"os"

	jsonextract "github.com/blainemoser/JsonExtract"
	"github.com/blainemoser/gitsync/utils"
)

type Configs struct {
	directories []string
}

// NewConfigs returns a new instance of Configs
func NewConfigs() *Configs {
	return &Configs{}
}

// SetDirectories sets the directories for the Configs instance
func (c *Configs) SetDirectories(path string, make string) (*Configs, error) {
	setup, err := os.OpenFile(path, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		if len(make) > 0 && os.IsNotExist(err) {
			err = utils.MakeDir(make)
			if err != nil {
				return c, err
			}
			c.directories = []string{make}
			return c, nil
		}
		return c, err
	}
	return c, c.parseDirectoryFile(setup)
}

// GetDirectories gets all directories
func (c *Configs) GetDirectories() []string {
	return c.directories
}

func (c *Configs) parseDirectoryFile(setup *os.File) error {
	stat, err := setup.Stat()
	if err != nil {
		return err
	}
	b := make([]byte, stat.Size())
	_, err = setup.Read(b)
	if err != nil {
		return err
	}
	return c.setRepos(b)
}

func (c *Configs) setRepos(b []byte) error {
	// Decode the contents
	json := jsonextract.JSONExtract{RawJSON: string(b)}
	directories, err := json.Extract("directories")
	if err != nil {
		return err
	}
	return c.populateRepos(directories)
}

func (c *Configs) populateRepos(directories interface{}) error {
	errors := make([]error, 0)
	paths := make([]string, 0)
	if repos, ok := directories.([]interface{}); ok {
		for _, repo := range repos {
			if directory, ok := repo.(string); ok {
				paths = append(paths, directory)
			} else {
				errors = append(errors, fmt.Errorf("%#v is not a path string", repo))
			}
		}
		if len(errors) > 0 {
			return utils.ParseErrors(errors)
		}
		c.directories = paths
		return nil
	}
	return fmt.Errorf("directories could not be parsed")
}
