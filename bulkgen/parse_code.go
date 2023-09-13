package bulkgen

import (
	"errors"
	"strconv"
	"strings"
)

func ParseAreaCode(code string) ([]int, error) {
	splittedCodes := strings.Split(code, ",")
	var trimmedCodes []string
	for _, code := range splittedCodes {
		trimCode := strings.TrimSpace(code)
		if len(trimCode) != 0 {
			trimmedCodes = append(trimmedCodes, trimCode)
		}
	}

	var parsedCodes []int
	for _, code := range trimmedCodes {
		numCode, err := strconv.Atoi(code)
		if err == nil {
			parsedCodes = append(parsedCodes, numCode)
		}
	}

	if len(parsedCodes) == 0 {
		err := errors.New("invalid input. Exiting Program")
		return nil, err
	}
	return parsedCodes, nil

}
