package util

import (
	"regexp"
	"strings"
)

func sanitizeFileName(fileName string) string {
	sanitized := strings.ReplaceAll(fileName, " ", "_")
	return regexp.MustCompile(`[^a-zA-Z0-9\._-]`).ReplaceAllString(sanitized, "")
}