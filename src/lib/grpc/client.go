package grpc

import (
	"context"
	"fmt"
	"main/log"
	"net"
	"os"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const socketPath = "/run/aikido.sock"

var conn *grpc.ClientConn
var client protos.AikidoClient
var rootCtx = context.Background()
var mainCtx, mainCxl = context.WithCancel(rootCtx)

func unixDialer(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	log.Debugf("Unix dialer called: %v %v", ctx, addr)
	return d.DialContext(ctx, "unix", addr)
}

func Init() {
	var err error
	conn, err = grpc.Dial(
		"unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(unixDialer),
	)
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}

	client = protos.NewAikidoClient(conn)

	log.Debug("Initialized gRPC client!")
}

func Uninit() {
	mainCxl()
	conn.Close()
}

func OnReceiveToken() {
	log.Infof("Client: %v", client)

	token := os.Getenv("AIKIDO_TOKEN")
	if token == "" {
		log.Warn("AIKIDO_TOKEN not found in env variables!")
		return
	}
	log.Info("Sending token: ", token)

	ctx, cancel := context.WithTimeout(mainCtx, 1*time.Second)
	defer cancel()

	_, err := client.OnReceiveToken(ctx, &protos.Token{Token: token})
	if err != nil {
		log.Debugf("Could not send token %v: %v", token, err)
	}
	log.Info("Token sent")
}

func OnReceiveLogLevel() {
	log.Infof("Client: %v", client)

	log_level := os.Getenv("AIKIDO_LOG_LEVEL")
	if log_level == "" {
		log.Debug("AIKIDO_LOG_LEVEL not found in env variables!")
		return
	}
	log.Info("Sending log level: ", log_level)

	ctx, cancel := context.WithTimeout(mainCtx, 1*time.Second)
	defer cancel()

	_, err := client.OnReceiveLogLevel(ctx, &protos.LogLevel{LogLevel: log_level})
	if err != nil {
		log.Debugf("Could not send log level %v: %v", log_level, err)
	}
	log.Info("Log level sent")
}

func OnReceiveDomain(domain string) {
	log.Infof("Client: %v", client)

	ctx, cancel := context.WithTimeout(mainCtx, 1*time.Second)
	defer cancel()

	log.Infof("Ctx: %v", ctx)

	_, err := client.OnReceiveDomain(ctx, &protos.Domain{Domain: domain})
	if err != nil {
		log.Debugf("Could not send domain %v: %v", domain, err)
	}
	log.Debugf("Domain sent: %v", domain)
}
