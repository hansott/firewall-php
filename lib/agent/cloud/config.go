package cloud

import (
	"encoding/json"
	. "main/aikido_types"
	"main/globals"
	"main/log"
)

func CheckConfigUpdatedAt() {
	response, err := SendCloudRequest(globals.EnvironmentConfig.ConfigEndpoint, globals.ConfigUpdatedAtAPI, globals.ConfigUpdatedAtMethod, nil)
	if err != nil {
		log.Debug("Error in sending polling config request: ", err)
		return
	}

	cloudConfigUpdatedAt := CloudConfigUpdatedAt{}
	err = json.Unmarshal(response, &cloudConfigUpdatedAt)
	if err != nil {
		return
	}

	if cloudConfigUpdatedAt.ConfigUpdatedAt == globals.CloudConfig.ConfigUpdatedAt {
		return
	}

	response, err = SendCloudRequest(globals.EnvironmentConfig.Endpoint, globals.ConfigAPI, globals.ConfigAPIMethod, nil)
	if err != nil {
		log.Debug("Error in sending config request: ", err)
		return
	}
	StoreCloudConfig(response)
}
