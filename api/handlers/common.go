package handlers

import (
	"strings"
)

func sanitizeURL(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	return "http://" + url
}
