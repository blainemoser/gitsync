package utils

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

// ParseErrors compacts a list of errors into a single error
func ParseErrors(errors []error) error {
	if len(errors) < 1 {
		return nil
	}
	result := []string{}
	for _, v := range errors {
		if v == nil {
			continue
		}
		result = append(result, v.Error())
	}
	if len(result) < 1 {
		return nil
	}
	return fmt.Errorf(strings.Join(result, "; "))
}

// MakeDir makes a directory
func MakeDir(path string) error {
	err := os.Mkdir(path, fs.FileMode(0777))
	if os.IsExist(err) {
		fmt.Printf("warning: %s\n", err.Error())
		return nil
	}
	return err
}

// GetPWD gets the present working directory
func GetPWD() (string, error) {
	result, err := exec.Command("pwd").CombinedOutput()
	return strings.ReplaceAll(string(result), "\n", ""), err
}
