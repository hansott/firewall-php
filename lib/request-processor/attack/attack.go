package attack

import (
	"main/ipc/protos"
	"main/utils"
)

func GetAttackDetectedMetadataProto(metadata map[string]string) []*protos.Metadata {
	var metadataProto []*protos.Metadata
	for key, value := range metadata {
		metadataProto = append(metadataProto, &protos.Metadata{Key: key, Value: value})
	}
	return metadataProto
}

func GetAttackDetectedProto(res utils.InterceptorResult) protos.AttackDetected {
	return protos.AttackDetected{
		Request: &protos.Request{},
		Attack: &protos.Attack{
			Kind:      string(res.Kind),
			Operation: res.Operation,
			Blocked:   utils.IsBlockingEnabled(),
			Source:    res.Source,
			Path:      res.PathToPayload,
			Payload:   res.Payload,
			Metadata:  GetAttackDetectedMetadataProto(res.Metadata),
			//UserId:    GetUserId(),
		},
	}
}
