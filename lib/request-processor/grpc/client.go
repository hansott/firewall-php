package grpc

import (
	"context"
	"fmt"
	"main/globals"
	"main/log"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const socketPath = "/run/aikido-" + globals.Version + ".sock"

var conn grpc.ClientConn
var client protos.AikidoClient

func Init() {
	conn, err := grpc.Dial(
		"unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}

	client = protos.NewAikidoClient(conn)

	log.Debugf("Current connection state: %s\n", conn.GetState().String())

	startCloudConfigRoutine()
}

func Uninit() {
	conn.Close()
}

/* Send outgoing domain to Aikido Agent via gRPC */
func OnDomain(domain string) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := client.OnDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Warnf("Could not send domain %v: %v", domain, err)
	}

	log.Debugf("Domain sent via socket: %v", domain)
}

/* Send request metadata (route & method) to Aikido Agent via gRPC */
func OnRequest(method string, route string, timeout time.Duration) bool {
	if client == nil {
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	requestStatus, err := client.OnRequest(ctx, &protos.RequestMetadata{Method: method, Route: route})
	if err != nil {
		log.Warnf("Could not send request metadata %v %v: %v", method, route, err)
		return true
	}

	log.Debugf("Request metadata sent via socket (%v %v) and got reply (%v)", method, route, requestStatus)
	return requestStatus.ForwardToServer
}

/* Get latest cloud config from Aikido Agent via gRPC */
func GetCloudConfig() {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cloudConfig, err := client.GetCloudConfig(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warnf("Could not get cloud config: %v", err)
	}

	log.Debugf("Got cloud config: %v", cloudConfig)
	setCloudConfig(cloudConfig)
}
