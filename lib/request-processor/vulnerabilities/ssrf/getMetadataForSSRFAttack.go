package ssrf

import "fmt"

func getMetadataForSSRFAttack(hostname string, port int) map[string]string {
	metadata := map[string]string{
		"hostname": hostname,
	}

	if port != 0 {
		metadata["port"] = fmt.Sprintf("%d", port)
	}

	return metadata
}
