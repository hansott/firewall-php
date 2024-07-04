package grpc

import (
	"context"
	"fmt"
	"main/log"
	"os"
	"sync"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const socketPath = "/run/aikido.sock"

var conn *grpc.ClientConn
var client protos.AikidoClient
var mu sync.Mutex

func Init() {
	mu.Lock()
	defer mu.Unlock()

	var err error
	conn, err = grpc.Dial(
		"unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}

	client = protos.NewAikidoClient(conn)

	log.Debug("Initialized gRPC client!")
}

func Uninit() {
	mu.Lock()
	defer mu.Unlock()

	conn.Close()
}

func OnReceiveToken() {
	mu.Lock()
	defer mu.Unlock()

	log.Infof("Client: %v", client)

	token := os.Getenv("AIKIDO_TOKEN")
	if token == "" {
		log.Warn("AIKIDO_TOKEN not found in env variables!")
		return
	}
	log.Info("Sending token: ", token)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.OnReceiveToken(ctx, &protos.Token{Token: token})
	if err != nil {
		log.Debugf("Could not send token %v: %v", token, err)
	}
	log.Info("Token sent")
}

func OnReceiveLogLevel() {
	mu.Lock()
	defer mu.Unlock()

	log.Infof("Client: %v", client)

	log_level := os.Getenv("AIKIDO_LOG_LEVEL")
	if log_level == "" {
		log.Debug("AIKIDO_LOG_LEVEL not found in env variables!")
		return
	}
	log.Info("Sending log level: ", log_level)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.OnReceiveLogLevel(ctx, &protos.LogLevel{LogLevel: log_level})
	if err != nil {
		log.Debugf("Could not send log level %v: %v", log_level, err)
	}
	log.Info("Log level sent")
}

func OnReceiveDomain(domain string) {
	mu.Lock()
	defer mu.Unlock()

	log.Infof("Client: %v", client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Setup gRPC logging to a file
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(log.Logger.Writer(), log.Logger.Writer(), log.Logger.Writer()))

	log.Infof("Client: %v", client)

	_, err := client.OnReceiveDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Debugf("Could not send domain %v: %v", domain, err)
	}
	log.Debugf("Domain sent: %v", domain)
}
