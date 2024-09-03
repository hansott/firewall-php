package cloud

import (
	. "main/aikido_types"
	"main/globals"
	"main/ipc/protos"
	"main/log"
	"main/utils"
)

func GetHeaders(grpcRequest *protos.Request) map[string][]string {
	headers := map[string][]string{}
	return headers
}

func GetRequestInfo(grpcRequest *protos.Request) RequestInfo {
	return RequestInfo{
		Method:    grpcRequest.Method,
		IPAddress: grpcRequest.IpAddress,
		UserAgent: grpcRequest.UserAgent,
		URL:       grpcRequest.Url,
		Headers:   GetHeaders(grpcRequest),
		Body:      grpcRequest.Body,
		Source:    grpcRequest.Source,
		Route:     grpcRequest.Route,
	}
}

func GetAttackDetails(grpcAttack *protos.Attack) AttackDetails {
	return AttackDetails{
		Kind:      grpcAttack.Kind,
		Operation: grpcAttack.Operation,
		Module:    grpcAttack.Module,
		Blocked:   grpcAttack.Blocked,
		Source:    grpcAttack.Source,
		Path:      grpcAttack.Path,
		Stack:     grpcAttack.Stack,
		Payload:   grpcAttack.Payload,
		//Metadata:  reqAttack.Metadata,
		User: utils.GetUserById(grpcAttack.UserId),
	}
}

func SendDetectedAttackEvent(req *protos.AttackDetected) {
	detectedAttackEvent := DetectedAttack{
		Type:    "detected_attack",
		Agent:   GetAgentInfo(),
		Request: GetRequestInfo(req.Request),
		Attack:  GetAttackDetails(req.Attack),
		Time:    utils.GetTime(),
	}

	response, err := SendCloudRequest(globals.EnvironmentConfig.Endpoint, globals.EventsAPI, globals.EventsAPIMethod, detectedAttackEvent)
	if err != nil {
		log.Warn("Error in sending detected attack event: ", err)
		return
	}
	StoreCloudConfig(response)
}
