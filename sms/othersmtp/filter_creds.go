package othersmtp

import (
	"errors"
	"strconv"
	"strings"
)

func FilterCreds(smtp string) ([]string, int, error) {
	trimCreds := strings.TrimSpace(smtp)
	splittedCreds := strings.Split(trimCreds, ",")
	if len(splittedCreds) != 4 {
		err := errors.New("invalid smtp format")
		return nil, 0, err
	}
	var filteredCreds []string
	for _, v := range splittedCreds {
		vTrimmed := strings.TrimSpace(v)
		if len(vTrimmed) != 0 {
			filteredCreds = append(filteredCreds, vTrimmed)
		}
	}
	if len(filteredCreds) != 4 {
		err := errors.New("invalid smtp format")
		return nil, 0, err
	}
	port, err := strconv.Atoi(filteredCreds[1])
	if err != nil {
		return nil, 0, err
	}
	return filteredCreds, port, nil
}
