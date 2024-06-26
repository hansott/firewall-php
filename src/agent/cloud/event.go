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

func SendEvent(route string, method string, payload interface{}) ([]byte, error) {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	apiEndpoint, err := url.JoinPath(globals.LocalConfig.Endpoint, route)
	if err != nil {
		return nil, fmt.Errorf("failed to build API endpoint: %v", err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	log.Debugf("Sending %s request to %s with content %s", method, apiEndpoint, jsonData)

	req, err := http.NewRequest(method, apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", globals.Token)
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

	log.Debug("Got response: ", string(body))

	return body, nil
}
