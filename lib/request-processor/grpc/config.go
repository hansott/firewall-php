package grpc

import (
	. "main/aikido_types"
	"main/globals"
	"main/ipc/protos"
	"time"
)

var (
	stopChan          chan struct{}
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
			AllowedIPAddresses: map[string]bool{},
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

	globals.CloudConfig.BypassedIps = map[string]bool{}
	for _, ip := range cloudConfigFromAgent.BypassedIps {
		globals.CloudConfig.BypassedIps[ip] = true
	}

	if cloudConfigFromAgent.Block {
		globals.CloudConfig.Block = 1
	} else {
		globals.CloudConfig.Block = 0
	}
}

func startCloudConfigRoutine() {
	GetCloudConfig()

	stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-cloudConfigTicker.C:
				GetCloudConfig()
			case <-stopChan:
				cloudConfigTicker.Stop()
				return
			}
		}
	}()
}

func stopCloudConfigRoutine() {
	if stopChan != nil {
		close(stopChan)
	}
}
