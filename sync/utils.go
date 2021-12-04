package sync

import (
	"fmt"
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
