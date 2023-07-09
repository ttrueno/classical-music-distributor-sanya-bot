package conv

import (
	"fmt"
	"strconv"
)

func Itoa64(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Atoi64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func ParseFirstInt64(str string) (int64, error) {
	startIndex := 0
	for i, char := range str {
		if (char >= '0' && char <= '9') || char == '-' || char == '+' {
			startIndex = i
		}
	}

	endIndex := startIndex
	for i, char := range str {
		if char >= '0' && char <= '9' {
			endIndex = i
		} else {
			break
		}
	}

	intStr := str[startIndex : endIndex+1]
	num, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int64: %w", err)
	}

	return num, nil
}
