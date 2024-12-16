package grpc

import (
	. "main/aikido_types"
	"main/globals"
	"main/ipc/protos"
	"main/log"
	"time"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

var (
	stopChan          chan struct{}
	cloudConfigTicker = time.NewTicker(1 * time.Minute)
)

func buildBlockedIpsTrie(name string, ipsList []string) (trieV4 *ipaddr.IPv4AddressTrie, trieV6 *ipaddr.IPv6AddressTrie) {
	if len(ipsList) == 0 {
		log.Debugf("%s -> Empty blocked IPs list!", name)
		return nil, nil
	} else {
		trieV4 = &ipaddr.IPv4AddressTrie{}
		trieV6 = &ipaddr.IPv6AddressTrie{}
		for _, ip := range ipsList {
			ipAddress, err := ipaddr.NewIPAddressString(ip).ToAddress()
			if err != nil {
				log.Infof("Invalid address: %s\n", ip)
				continue
			}

			if ipAddress.IsIPv4() {
				trieV4.Add(ipAddress.ToIPv4())
			} else if ipAddress.IsIPv6() {
				trieV6.Add(ipAddress.ToIPv6())
			}
		}
	}

	log.Debugf("%s (v4): %v", name, trieV4)
	log.Debugf("%s (v6): %v", name, trieV6)
	return trieV4, trieV6
}

func setCloudConfig(cloudConfigFromAgent *protos.CloudConfig) {
	if cloudConfigFromAgent == nil {
		return
	}

	globals.CloudConfigMutex.Lock()
	defer globals.CloudConfigMutex.Unlock()

	globals.CloudConfig.ConfigUpdatedAt = cloudConfigFromAgent.ConfigUpdatedAt

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

	globals.CloudConfig.GeoBlockedIpsTrieV4, globals.CloudConfig.GeoBlockedIpsTrieV6 = buildBlockedIpsTrie("geoip", cloudConfigFromAgent.GeoBlockedIps)
	globals.CloudConfig.TorBlockedIpsTrieV4, globals.CloudConfig.TorBlockedIpsTrieV6 = buildBlockedIpsTrie("tor", cloudConfigFromAgent.TorBlockedIps)
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
