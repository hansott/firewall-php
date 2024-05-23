package main

import (
	"log"
	"net/url"
)

func ExitIfKeyDoesNotExistInMap[K comparable, V any](m map[K]V, key K) {
	if _, exists := m[key]; !exists {
		log.Fatalf("Key %s does not exist in map!", key)
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
		log.Fatalf("Error parsing JSON: key %s has incorrect type", key)
	}
	return *value
}

func GetDomain(rawurl string) string {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}
