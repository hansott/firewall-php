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
)

type server struct {
	protos.AikidoServer
}

func (s *server) OnReceiveDomain(ctx context.Context, req *protos.Domain) (*protos.BoolResponse, error) {
	log.Debugf("Received domain: %s", req.GetDomain())
	globals.HostnamesMutex.Lock()
	defer globals.HostnamesMutex.Unlock()

	globals.Hostnames[req.GetDomain()] = true
	return &protos.BoolResponse{Success: true}, nil
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

func Uninit() {
	// Remove the socket file if it exists
	if _, err := os.Stat(globals.SocketPath); err == nil {
		if err := os.RemoveAll(globals.SocketPath); err != nil {
			panic(fmt.Sprintf("failed to remove existing socket: %v", err))
		}
	}
}
