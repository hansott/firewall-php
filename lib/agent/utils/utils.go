package utils

import (
	"fmt"
	"time"
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

func StartPollingRoutine(stopChan chan struct{}, ticker *time.Ticker, pollingFunction func()) {
	go func() {
		for {
			select {
			case <-ticker.C:
				pollingFunction()
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func StopPollingRouting(stopChan chan struct{}) {
	close(stopChan)
}
