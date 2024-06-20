package main

import (
	"context"
	"fmt"
	"main/log"
	"time"

	"main/ipc/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const socketPath = "/var/aikido.sock"

var conn *grpc.ClientConn

func InitClient() {
	conn, err := grpc.Dial(
		"unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}

	c := protos.NewAikidoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SendDomain(ctx, &protos.Domain{Domain: "www.example.com"})
	if err != nil {
		panic(fmt.Sprintf("could not greet: %v", err))
	}

	log.Info("gRPC client initialized!")
}

func UninitClient() {
	conn.Close()
}
