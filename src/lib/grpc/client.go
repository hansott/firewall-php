package grpc

import (
	"context"
	"fmt"
	"main/log"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
)

const socketPath = "/var/aikido.sock"

var conn *grpc.ClientConn
var cancel context.CancelFunc
var client protos.AikidoClient
var ctx context.Context

func Init() {
	var err error
	conn, err = grpc.Dial(
		"unix://"+socketPath,
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}

	client = protos.NewAikidoClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	log.Debug("Initialized gRPC client!")
}

func Uninit() {
	conn.Close()
	cancel()
}

func SendDomain(domain string) {
	_, err := client.SendDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Debugf("Could not send domain: %v", domain, err)
	}
}
