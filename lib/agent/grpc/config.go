package grpc

import (
	"main/cloud"
	"main/globals"
	"main/log"
)

func storeConfig(token, logLevel string, blocking, localhostAllowedByDefault, collectApiSchema bool) {
	globals.AikidoConfig.ConfigMutex.Lock()
	defer globals.AikidoConfig.ConfigMutex.Unlock()

	if globals.AikidoConfig.Token != "" || token == "" {
		// Token is already set or the newly provided token via gRPC is invalid -> do not update the config in this case
		log.Debugf("Aikido Config is already set. Returning early.")
		return
	}

	globals.AikidoConfig.Token = token
	globals.AikidoConfig.LogLevel = logLevel
	globals.AikidoConfig.Blocking = blocking
	globals.AikidoConfig.LocalhostAllowedByDefault = localhostAllowedByDefault
	globals.AikidoConfig.CollectApiSchema = collectApiSchema

	log.SetLogLevel(logLevel)
	cloud.SendStartEvent()

	log.Debugf("Updated Aikido Config with the one passed via gRPC!")
}
