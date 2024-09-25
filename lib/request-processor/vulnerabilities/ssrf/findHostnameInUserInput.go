package ssrf

import (
	"main/helpers"
)

func findHostnameInUserInput(userInput string, hostname string, port int) bool {

	if len(userInput) <= 1 {
		return false
	}

	hostnameURL := helpers.TryParseURL("http://" + hostname)
	if hostnameURL == nil {
		return false
	}

	variants := []string{userInput, "http://" + userInput, "https://" + userInput}

	for _, variant := range variants {
		userInputURL := helpers.TryParseURL(variant)
		if userInputURL != nil && userInputURL.Hostname() == hostnameURL.Hostname() {
			userPort := helpers.GetPortFromURL(userInputURL)

			if port == 0 {
				return true
			}

			if userPort == port {
				return true
			}
		}
	}

	return false
}
