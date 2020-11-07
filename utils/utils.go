package utils

import (
	"strings"
)

func ShouldCountFile(file string, ignoreFiles []string) bool {
	if !strings.HasSuffix(file, ".go") {
		return false
	}

	for _, suffix := range ignoreFiles {
		if strings.HasSuffix(file, suffix) {
			return false
		}
	}

	return true
}
