package ssrf

import (
	"reflect"
	"testing"
)

func TestGetMetadataForSSRFAttack(t *testing.T) {
	// Test case: port is undefined (equivalent in Go is port 0)
	t.Run("port is undefined", func(t *testing.T) {
		expected := map[string]string{
			"hostname": "example.com",
		}
		result := getMetadataForSSRFAttack("example.com", 0)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case: port is defined
	t.Run("port is defined", func(t *testing.T) {
		expected := map[string]string{
			"hostname": "example.com",
			"port":     "80",
		}
		result := getMetadataForSSRFAttack("example.com", 80)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	// Test case: port is 443
	t.Run("port is 443", func(t *testing.T) {
		expected := map[string]string{
			"hostname": "example.com",
			"port":     "443",
		}
		result := getMetadataForSSRFAttack("example.com", 443)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})
}
