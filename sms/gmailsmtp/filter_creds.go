package gmailsmtp

import (
	"errors"
	"strings"
)

func FilterCreds(smtp string) ([]string, error) {
	splittedCreds := strings.Split(smtp, ":")
	if len(splittedCreds) != 2 {
		err := errors.New("invalid smtp format")
		return nil, err
	}
	return splittedCreds, nil
}
