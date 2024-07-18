package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func CheckIfKeyExists[K comparable, V any](m map[K]V, key K) {
	if _, exists := m[key]; !exists {
		panic(fmt.Sprintf("Key %v does not exist in map!", key))
	}
}

func GetFromMap[T any](m map[string]interface{}, key string) *T {
	value, ok := m[key]
	if !ok {
		return nil
	}
	result, ok := value.(T)
	if !ok {
		return nil
	}
	return &result
}

func MustGetFromMap[T any](m map[string]interface{}, key string) T {
	value := GetFromMap[T](m, key)
	if value == nil {
		panic(fmt.Sprintf("Error parsing JSON: key %s does not exist or it has an incorrect type", key))
	}
	return *value
}

func FixURL(url string) string {
	if !strings.HasPrefix(url, "https://") && strings.HasPrefix(url, "https:/") {
		return strings.Replace(url, "https:/", "https://", 1)
	}
	if !strings.HasPrefix(url, "http://") && strings.HasPrefix(url, "http:/") {
		return strings.Replace(url, "http:/", "http://", 1)
	}
	return url
}

func GetDomain(rawurl string) string {
	parsedURL, err := url.Parse(FixURL(rawurl))
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}
