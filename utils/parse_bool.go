package utils

import (
	"strconv"
	"strings"
)

func AsBool(ptr *string) bool {
	if ptr == nil {
		return false
	}
	b, _ := strconv.ParseBool(strings.TrimSpace(*ptr))
	return b
}
