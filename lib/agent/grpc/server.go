package grpc

import (
	"context"
	"fmt"
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

func (s *server) OnDomain(ctx context.Context, req *protos.Domain) (*emptypb.Empty, error) {
	log.Debugf("Received domain: %s:%d", req.GetDomain(), req.GetPort())
	storeDomain(req)
	return &emptypb.Empty{}, nil
}

func (s *server) OnRequestInit(ctx context.Context, req *protos.RequestMetadataInit) (*protos.RequestStatus, error) {
	log.Debugf("Received request metadata: %s %s", req.GetMethod(), req.GetRoute())

	go storeStats()

	return getRequestStatus(req), nil
}

func (s *server) OnRequestShutdown(ctx context.Context, req *protos.RequestMetadataShutdown) (*emptypb.Empty, error) {
	log.Debugf("Received request metadata: %s %s %d", req.GetMethod(), req.GetRoute(), req.GetStatusCode())

	go storeRoute(req)
	go updateRateLimitingStatus(req)

	return &emptypb.Empty{}, nil
}

func (s *server) GetCloudConfig(ctx context.Context, req *emptypb.Empty) (*protos.CloudConfig, error) {
	return getCloudConfig(), nil
}

func (s *server) OnUser(ctx context.Context, req *protos.User) (*emptypb.Empty, error) {
	log.Debugf("Received user event: %s", req.GetId())
	go onUserEvent(req.GetId(), req.GetUsername(), req.GetIp())
	return &emptypb.Empty{}, nil
}

func StartServer(lis net.Listener) {
	s := grpc.NewServer()
	protos.RegisterAikidoServer(s, &server{})

	log.Infof("Server is running on Unix socket %s", globals.SocketPath)
	if err := s.Serve(lis); err != nil {
		log.Warnf("gRPC server failed to serve: %v", err)
	}
	log.Warnf("gRPC server went down!")
	lis.Close()
}

func Init() bool {
	// Remove the socket file if it already exists
	if _, err := os.Stat(globals.SocketPath); err == nil {
		os.RemoveAll(globals.SocketPath)
	}

	lis, err := net.Listen("unix", globals.SocketPath)
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	// Change the permissions of the socket to make it accessible by non-root users
	if err := os.Chmod(globals.SocketPath, 0777); err != nil {
		panic(fmt.Sprintf("failed to change permissions of Unix socket: %v", err))
	}

	go StartServer(lis)
	return true
}

func Uninit() {
	// Remove the socket file if it exists
	if _, err := os.Stat(globals.SocketPath); err == nil {
		if err := os.RemoveAll(globals.SocketPath); err != nil {
			panic(fmt.Sprintf("failed to remove existing socket: %v", err))
		}
	}
}
