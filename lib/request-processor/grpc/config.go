package grpc

import (
	"errors"
	. "main/aikido_types"
	"main/globals"
	"main/ipc/protos"
	"time"
)

var (
	stop              chan struct{}
	cloudConfigTicker = time.NewTicker(1 * time.Minute)
)

func setCloudConfig(cloudConfigFromAgent *protos.CloudConfig) {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	globals.CloudConfig.Endpoints = map[EndpointKey]EndpointData{}
	for _, ep := range cloudConfigFromAgent.Endpoints {
		endpointData := EndpointData{
			ForceProtectionOff: ep.ForceProtectionOff,
			RateLimiting: RateLimiting{
				Enabled: ep.RateLimiting.Enabled,
			},
		}
		for _, ip := range ep.AllowedIPAddresses {
			endpointData.AllowedIPAddresses[ip] = true
		}
		globals.CloudConfig.Endpoints[EndpointKey{Method: ep.Method, Route: ep.Route}] = endpointData
	}

	globals.CloudConfig.BlockedUserIds = map[string]bool{}
	for _, userId := range cloudConfigFromAgent.BlockedUserIds {
		globals.CloudConfig.BlockedUserIds[userId] = true
	}

	globals.CloudConfig.IpsExcludedFromRateLimiting = map[string]bool{}
	for _, ip := range cloudConfigFromAgent.IpsExcludedFromRateLimiting {
		globals.CloudConfig.IpsExcludedFromRateLimiting[ip] = true
	}
}

func GetEndpointConfig(method string, route string) (EndpointData, error) {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	endpointData, exists := globals.CloudConfig.Endpoints[EndpointKey{Method: method, Route: route}]
	if !exists {
		return EndpointData{}, errors.New("endpoint does not exist")
	}

	return endpointData, nil
}

func startCloudConfigRoutine() {
	GetCloudConfig()

	stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-cloudConfigTicker.C:
				GetCloudConfig()
			case <-stop:
				cloudConfigTicker.Stop()
				return
			}
		}
	}()
}
