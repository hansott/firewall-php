package main

import (
	"fmt"
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

	log.Init()
	config.Init()
	machine.Init()
	go grpc.Init()

	runUntilKilled()

	grpc.Uninit()
	log.Uninit()
	log.Infof("Aikido agent v%s stopped!", globals.Version)
}
