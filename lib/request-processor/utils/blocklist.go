package utils

import (
	"main/aikido_types"
	"main/log"
	"net/netip"

	"go4.org/netipx"
)

func BuildIpBlocklist(name, description string, ipsList []string) (*aikido_types.IpBlockList, error) {
	trieBuilder := netipx.IPSetBuilder{}

	for _, ip := range ipsList {
		prefix, err := netip.ParsePrefix(ip)
		if err == nil {
			trieBuilder.AddPrefix(prefix)
		} else {
			parsedIP, err := netip.ParseAddr(ip)
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
