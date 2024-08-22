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
func OnDomain(domain string, port int) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := client.OnDomain(ctx, &protos.Domain{Domain: domain, Port: int32(port)})
	if err != nil {
		log.Warnf("Could not send domain %v: %v", domain, err)
		return
	}

	log.Debugf("Domain sent via socket: %v:%v", domain, port)
}

/* Send request metadata (route & method) to Aikido Agent via gRPC */
func OnRequestInit(method string, route string, timeout time.Duration) bool {
	if client == nil {
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	requestStatus, err := client.OnRequestInit(ctx, &protos.RequestMetadataInit{Method: method, Route: route})
	if err != nil {
		log.Warnf("Could not send request metadata %v %v: %v", method, route, err)
		return true
	}

	log.Debugf("Request metadata sent via socket (%v %v) and got reply (%v)", method, route, requestStatus)
	return requestStatus.ForwardToServer
}

/* Send request metadata (route, method & status code) to Aikido Agent via gRPC */
func OnRequestShutdown(method string, route string, statusCode int, timeout time.Duration) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := client.OnRequestShutdown(ctx, &protos.RequestMetadataShutdown{Method: method, Route: route, StatusCode: int32(statusCode)})
	if err != nil {
		log.Warnf("Could not send request metadata %v %v %v: %v", method, route, statusCode, err)
		return
	}

	log.Debugf("Request metadata sent via socket (%v %v %v)", method, route, statusCode)
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
		return
	}

	log.Debugf("Got cloud config: %v", cloudConfig)
	setCloudConfig(cloudConfig)
}
