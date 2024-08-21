package cloud

import (
	. "main/aikido_types"
	"main/globals"
	"main/log"
)

func SendStartEvent() {
	startedEvent := Started{
		Type:  "started",
		Agent: GetAgentInfo(),
		Time:  GetTime(),
	}

	response, err := SendCloudRequest(globals.EnvironmentConfig.Endpoint, globals.EventsAPI, globals.EventsAPIMethod, startedEvent)
	if err != nil {
		log.Warn("Error in sending start event: ", err)
		return
	}
	StoreCloudConfig(response)
}
