package mserv

import (
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	addr   string
	server *grpc.Server
}

func NewGRPCServer(addr string, server *grpc.Server) (Server, error) {
	if addr == "" {
		return nil, errors.New("missing bind address for grpc.Server")
	}

	return &GRPCServer{
		addr:   addr,
		server: server,
	}, nil
}

func (g *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		if err := g.server.Serve(lis); err != nil {
			if err != grpc.ErrServerStopped {
				log.Fatalf("start grpc server %s error %s", g.addr, err)
			}
		}
	}()

	return nil
}

func (g *GRPCServer) Stop() error {
	if g.server == nil {
		return nil
	}

	g.server.GracefulStop()

	return nil
}
