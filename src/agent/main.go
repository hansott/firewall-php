package main

import (
	"fmt"
	"main/cloud"
	"main/config"
	"main/globals"
	"main/grpc"
	"main/log"
	"main/machine"
	"os"
	"os/signal"
	"syscall"
)

func runUntilKilled() {
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigChannel
		fmt.Printf("Received signal: %s\n", sig)
		done <- true
	}()

	<-done
}

func main() {
	log.SetLogLevel("INFO")
	log.Infof("Aikido agent v%s started!", globals.Version)

	config.Init()
	if err := log.SetLogLevel(globals.Config.LogLevel); err != nil {
		panic(fmt.Sprintf("Error setting log level: %s", err))
	}

	machine.Init()
	go cloud.Init()
	go grpc.Init()

	runUntilKilled()

	log.Infof("Aikido agent v%s stopped!", globals.Version)
}
