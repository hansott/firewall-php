package main

import "C"
import "net/url"

//export AikidoLibNormalizeDomain
func AikidoLibNormalizeDomain(rawurl string) string {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

func main() {}
