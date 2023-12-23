package gslib

import (
	"github.com/google/uuid"
	"regexp"
)

func IsValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^(?:http|ftp)s?://(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+(?:[a-zA-Z]{2,6}\.?|[a-zA-Z0-9-]{2,}\.?)|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|\[?[A-F0-9]*:[A-F0-9:]+\]?(?::\d+)?(?:/?|[/?]\S+)$`)

	return urlRegex.MatchString(url)
}

func GenerateID() string {
	return uuid.New().String()
}