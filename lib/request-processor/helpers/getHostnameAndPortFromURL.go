package helpers

func GetHostnameAndPortFromURL(url string) (string, int) {
	parsedURL := TryParseURL(url)
	if parsedURL == nil {
		return "", 0
	}
	return parsedURL.Hostname(), GetPortFromURL(parsedURL)
}
