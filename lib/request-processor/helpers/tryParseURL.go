package helpers

import (
	"net/url"
)

func TryParseURL(input string) *url.URL {
	parsedURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil
	}
	return parsedURL
}
