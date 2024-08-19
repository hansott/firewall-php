package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/globals"
	"main/log"
	"net/http"
	"net/url"
)

func SendCloudRequest(endpoint string, route string, method string, payload interface{}) ([]byte, error) {
	apiEndpoint, err := url.JoinPath(endpoint, route)
	if err != nil {
		return nil, fmt.Errorf("failed to build API endpoint: %v", err)
	}

	var buffer *bytes.Buffer = nil
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", err)
		}

		log.Infof("Sending %s request to %s with content %s", method, apiEndpoint, jsonData)

		buffer = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiEndpoint, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", globals.EnvironmentConfig.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	log.Info("Got response: ", string(body))
	return body, nil
}
