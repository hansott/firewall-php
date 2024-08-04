package grpc

import (
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
		globals.CloudConfig.Endpoints[EndpointKey{Method: ep.Method, Route: ep.Route}] = endpointData
	}

	for _, userId := range cloudConfigFromAgent.BlockedUserIds {
		globals.CloudConfig.BlockedUserIds[userId] = true
	}

	for _, allowedIpAddress := range cloudConfigFromAgent.AllowedIPAddresses {
		globals.CloudConfig.AllowedIPAddresses[allowedIpAddress] = true
	}
}

func IsRequestConfiguredForRateLimiting(method string, route string) bool {
	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	_, exists := globals.CloudConfig.Endpoints[EndpointKey{Method: method, Route: route}]

	return exists
}

func startCloudConfigRoutine() {
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
