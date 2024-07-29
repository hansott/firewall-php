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
}

func Uninit() {
	conn.Close()
}

func OnReceiveDomain(domain string) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := client.OnReceiveDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Warnf("Could not send domain %v: %v", domain, err)
	}

	log.Debugf("Domain sent via socket: %v", domain)
}

func OnReceiveHttpRequestInfo(method string, route string) {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.OnReceiveHttpRequestInfo(ctx, &protos.HttpRequestInfo{Method: method, Route: route})
	if err != nil {
		log.Warnf("Could not send http request info %v %v: %v", method, route, err)
	}

	log.Infof("Http request info sent via socket: %v %v", method, route)
}
