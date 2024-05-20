package main

import "C"
import "net/url"

//export NormalizeDomain
func NormalizeDomain(rawurl string) string {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

func main() {}
