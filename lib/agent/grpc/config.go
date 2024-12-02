package grpc

import (
	"main/globals"
	"main/log"
)

func storeConfig(token, logLevel string, blocking, localhostAllowedByDefault, collectApiSchema bool) {
	globals.AikidoConfig.ConfigMutex.Lock()
	defer globals.AikidoConfig.ConfigMutex.Unlock()

	if globals.AikidoConfig.Token != "" || token == "" {
		// Token is already set or the newly provided token via gRPC is invalid -> do not update the config in this case
		return
	}

	globals.AikidoConfig.Token = token
	globals.AikidoConfig.LogLevel = logLevel
	globals.AikidoConfig.Blocking = blocking
	globals.AikidoConfig.LocalhostAllowedByDefault = localhostAllowedByDefault
	globals.AikidoConfig.CollectApiSchema = collectApiSchema

	log.SetLogLevel(logLevel)
}
