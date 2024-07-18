package grpc

import (
	"context"
	"fmt"
	"main/log"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const socketPath = "/run/aikido.sock"

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

	log.Debugf("Client: %v", client)
	log.Debugf("Current connection state: %s\n", conn.GetState().String())
}

func Uninit() {
	conn.Close()
}

func OnReceiveDomain(domain string) {
	log.Debug("OnReceiveDomain")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.OnReceiveDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Debugf("Could not send domain %v: %v", domain, err)
	}

	log.Infof("Domain sent via socket: %v", domain)
}
