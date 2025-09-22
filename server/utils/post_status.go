package utils

import (
	"slices"
)

func IsValidStatus(status string) bool {
	validStatus := []string{
		"published",
		"draft",
		"deleted",
	}
	return slices.Contains(validStatus, status)
}
