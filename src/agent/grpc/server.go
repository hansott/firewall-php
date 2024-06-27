package grpc

import (
	"context"
	"fmt"
	"main/cloud"
	"main/globals"
	"main/ipc/protos"
	"main/log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	protos.AikidoServer
}

func (s *server) OnReceiveToken(ctx context.Context, req *protos.Token) (*emptypb.Empty, error) {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	newToken := req.GetToken()

	if globals.Token == newToken {
		// Got the same token, nothing to do
		return &emptypb.Empty{}, nil
	}

	log.Infof("Received new token: %s", newToken)

	if globals.Token != "" {
		// Token was previously set and got a new diffent token (token update)
		// Stop the previous cloud communication routines
		log.Infof("Stopping cloud routines for previous token...")
		cloud.Uninit()
	}

	globals.Token = newToken
	log.Infof("Starting cloud routines for new token...")
	go cloud.Init()

	return &emptypb.Empty{}, nil
}

func (s *server) OnReceiveLogLevel(ctx context.Context, req *protos.LogLevel) (*emptypb.Empty, error) {
	globals.ConfigMutex.Lock()
	defer globals.ConfigMutex.Unlock()

	newLogLevel := req.GetLogLevel()

	if globals.LogLevel == newLogLevel {
		// Got the same log level, nothing to do
		return &emptypb.Empty{}, nil
	}

	log.Infof("Received new log level: %s", newLogLevel)

	globals.LogLevel = newLogLevel
	log.SetLogLevel(globals.LogLevel)

	return &emptypb.Empty{}, nil
}

func (s *server) OnReceiveDomain(ctx context.Context, req *protos.Domain) (*emptypb.Empty, error) {
	log.Debugf("Received domain: %s", req.GetDomain())
	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	globals.Hostnames[req.GetDomain()] = true
	return &emptypb.Empty{}, nil
}

func Init() {
	// Remove the socket file if it already exists
	if _, err := os.Stat(globals.SocketPath); err == nil {
		if err := os.RemoveAll(globals.SocketPath); err != nil {
			panic(fmt.Sprintf("failed to remove existing socket: %v", err))
		}
	}

	lis, err := net.Listen("unix", globals.SocketPath)
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	defer lis.Close()

	// Change the permissions of the socket to make it accessible by non-root users
	if err := os.Chmod(globals.SocketPath, 0777); err != nil {
		panic(fmt.Sprintf("failed to change permissions of Unix socket: %v", err))
	}

	s := grpc.NewServer()
	protos.RegisterAikidoServer(s, &server{})

	log.Infof("Server is running on Unix socket %s", globals.SocketPath)
	if err := s.Serve(lis); err != nil {
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}
