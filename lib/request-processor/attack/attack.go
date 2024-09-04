package attack

import (
	"fmt"
	"main/context"
	"main/ipc/protos"
	"main/utils"
)

/* Convert metadata map to protobuf structure to be sent via gRPC to the Agent */
func GetMetadataProto(metadata map[string]string) []*protos.Metadata {
	var metadataProto []*protos.Metadata
	for key, value := range metadata {
		metadataProto = append(metadataProto, &protos.Metadata{Key: key, Value: value})
	}
	return metadataProto
}

/* Convert headers map to protobuf structure to be sent via gRPC to the Agent */
func GetHeadersProto() []*protos.Header {
	var headersProto []*protos.Header
	for key, value := range context.GetHeaders() {
		valueStr, ok := value.(string)
		if ok {
			headersProto = append(headersProto, &protos.Header{Key: key, Value: valueStr})
		}
	}
	return headersProto
}

/* Construct the AttackDetected protobuf structure to be sent via gRPC to the Agent */
func GetAttackDetectedProto(res utils.InterceptorResult) protos.AttackDetected {
	return protos.AttackDetected{
		Request: &protos.Request{
			Method:    context.GetMethod(),
			IpAddress: context.GetIp(),
			UserAgent: context.GetUserAgent(),
			Url:       context.GetUrl(),
			Headers:   GetHeadersProto(),
			Body:      context.GetBody(),
			Source:    "php",
			Route:     context.GetRoute(),
		},
		Attack: &protos.Attack{
			Kind:      string(res.Kind),
			Operation: res.Operation,
			Blocked:   utils.IsBlockingEnabled(),
			Source:    res.Source,
			Path:      res.PathToPayload,
			Payload:   res.Payload,
			Metadata:  GetMetadataProto(res.Metadata),
			UserId:    context.GetUserId(),
		},
	}
}

func BuildAttackDetectedMessage(result utils.InterceptorResult) string {
	return fmt.Sprintf("Aikido firewall has blocked %s: %s(...) originating from %s%s",
		utils.GetDisplayNameForAttackKind(result.Kind),
		result.Operation,
		result.Source,
		utils.EscapeHTML(result.PathToPayload))
}

func GetAttackDetectedAction(result utils.InterceptorResult) string {
	return fmt.Sprintf(`{"action": "throw", "message": "%s", "code": -1}`, BuildAttackDetectedMessage(result))
}
