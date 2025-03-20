package utils

import (
	"main/aikido_types"
	"main/log"

	"inet.af/netaddr"
)

func BuildIpBlocklist(name, description string, ipsList []string) (*aikido_types.IpBlockList, error) {
	trieBuilder := netaddr.IPSetBuilder{}

	for _, ip := range ipsList {
		prefix, err := netaddr.ParseIPPrefix(ip)
		if err == nil {
			trieBuilder.AddPrefix(prefix)
		} else {
			parsedIP, err := netaddr.ParseIP(ip)
			if err == nil {
				trieBuilder.Add(parsedIP)
			} else {
				log.Infof("Invalid address for %s: %s\n", name, ip)
				continue // Skip invalid IPs
			}
		}
	}

	trie, err := trieBuilder.IPSet()
	if err != nil {
		return nil, err
	}

	return &aikido_types.IpBlockList{
		Description: description,
		IpSet:       *trie,
	}, nil
}
