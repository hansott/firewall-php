package attack

import (
	"main/context"
	"main/ipc/protos"
	"main/utils"
)

func GetMetadataProto(metadata map[string]string) []*protos.Metadata {
	var metadataProto []*protos.Metadata
	for key, value := range metadata {
		metadataProto = append(metadataProto, &protos.Metadata{Key: key, Value: value})
	}
	return metadataProto
}

func GetHeadersProto() []*protos.Header {

	var headersProto []*protos.Header
	for key, value := range context.GetHeaders() {
		headersProto = append(headersProto, &protos.Header{Key: key, Value: value})
	}
	return headersProto
}

func GetAttackDetectedProto(res utils.InterceptorResult) protos.AttackDetected {
	return protos.AttackDetected{
		Request: &protos.Request{
			Method:    context.GetMethod(),
			IpAddress: context.GetIp(),
			UserAgent: context.GetUserAgent(),
			Url:       context.GetUrl(),
			Headers:   GetHeadersProto(),
			Body:      context.GetBodyRaw(),
			// source = 7;
			Route: context.GetRoute(),
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
