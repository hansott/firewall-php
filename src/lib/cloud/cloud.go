package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	. "main/aikido_types"
	. "main/globals"
	"main/log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	stop chan struct{}
	wg   sync.WaitGroup
)

func SendEvent(route string, method string, payload interface{}) (map[string]interface{}, error) {
	apiEndpoint, err := url.JoinPath(InitData.Aikido.Endpoint, route)
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

	req.Header.Set("Authorization", InitData.Aikido.Token)
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

	var responseJson map[string]interface{}
	err = json.Unmarshal(body, &responseJson)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return responseJson, nil
}

func SendStartEvent() {
	startedEvent := Started{
		Type: "started",
		Agent: AgentInfo{
			DryMode:   !InitData.Aikido.Blocking,
			Hostname:  InitData.Machine.HostName,
			Version:   InitData.Aikido.Version,
			IPAddress: InitData.Machine.IPAddress,
			OS: OsInfo{
				Name:    InitData.Machine.OS,
				Version: InitData.Machine.OSVersion,
			},
			Packages: make(map[string]string, 0),
			NodeEnv:  "",
		},
		Time: time.Now().Unix(),
	}

	_, err := SendEvent(EventsAPI, EventsAPIMethod, startedEvent)
	if err != nil {
		log.Debug("Error in sending start event: ", err)
	}
}

func StartPollingThread() {
	stop = make(chan struct{})

	ticker := time.NewTicker(time.Second * 10)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				return
				//MakeGetRequest(endpoint)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func StopPollingThread() {
	close(stop)
	wg.Wait()
}

func Init() {
	SendStartEvent()
	//StartPollingThread()
}

func Uninit() {
	//StopPollingThread()
}
