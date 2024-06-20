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

	response, err := SendEvent(globals.EventsAPI, globals.EventsAPIMethod, startedEvent)
	if err != nil {
		log.Debug("Error in sending start event: ", err)
	}
	UpdateCloudConfig(response)
}
