package main

import "C"
import (
	"fmt"
	"net/url"
)

func NormalizeDomain(rawurl string) string {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

//export OnEvent
func OnEvent(event string) string {
	fmt.Println("[AIKIDO-GO] OnEvent: ", event)
	return "{}"
}

func main() {}
