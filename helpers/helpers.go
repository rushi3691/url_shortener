package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	link := strings.Replace(url, "http://", "", 1)
	link = strings.Replace(link, "https://", "", 1)
	link = strings.Replace(link, "www.", "", 1)
	link = strings.Split(link, "/")[0]
	return link != os.Getenv("DOMAIN")
}
