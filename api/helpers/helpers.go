package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {

	if url[:4] != "http" {
		url = "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]
	if url == os.Getenv("DOMAIN") {
		return false
	}
	return true
}
