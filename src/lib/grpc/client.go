package grpc

import (
	"context"
	"fmt"
	"main/log"
	"os"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func SendToken() {
	token := os.Getenv("AIKIDO_TOKEN")
	if token == "" {
		log.Warn("AIKIDO_TOKEN not found in env variables!")
		return
	}
	log.Info("Sending token: ", token)
	_, err := client.SendToken(ctx, &protos.Token{Token: token})
	if err != nil {
		log.Debugf("Could not send token %v: %v", token, err)
	}
}

func SendDomain(domain string) {
	log.Info("Got domain: ", domain)
	_, err := client.SendDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Debugf("Could not send domain %v: %v", domain, err)
	}
}
