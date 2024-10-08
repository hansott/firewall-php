package utils

import (
	"errors"
	. "main/aikido_types"
	"main/globals"
)

func GetEndpointConfig(method string, route string) (EndpointData, error) {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	endpointData, exists := globals.CloudConfig.Endpoints[EndpointKey{Method: method, Route: route}]
	if !exists {
		return EndpointData{}, errors.New("endpoint does not exist")
	}

	return endpointData, nil
}

func AreEndpointsConfigured() bool {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	return len(globals.CloudConfig.Endpoints) != 0
}
