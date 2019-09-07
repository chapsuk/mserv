package mserv

import (
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// GRPCServerOption allows redefine default grpc server settings
type GRPCServerOption func(*GRPCServer)

// GRPCSkipErrors sets silent mode for Start/Stop grpc server, skip all errors
func GRPCSkipErrors(skip bool) GRPCServerOption {
	return func(s *GRPCServer) {
		s.skipErrors = skip
	}
}

// GRPCServer controll grpc server start/stop process
type GRPCServer struct {
	addr       string
	skipErrors bool
	server     *grpc.Server
}

// NewGRPCServer returns grpc server wrapper
func NewGRPCServer(addr string, server *grpc.Server, opts ...GRPCServerOption) Server {
	srv := &GRPCServer{
		addr:       addr,
		server:     server,
		skipErrors: false,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

// Start starts grpc server
func (g *GRPCServer) Start() error {
	if len(g.addr) == 0 {
		return g.returnErr(errors.New("missing bind addr"))
	}

	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return g.returnErr(err)
	}

	go func() {
		if err := g.server.Serve(lis); err != nil {
			if err != grpc.ErrServerStopped && !g.skipErrors {
				panic(fmt.Sprintf("grpc: Serve return unexpected error: %s", err))
			}
		}
	}()

	return nil
}

// Stop gracefully stops grpc server
func (g *GRPCServer) Stop() error {
	g.server.GracefulStop()
	return nil
}

func (g *GRPCServer) returnErr(err error) error {
	if g.skipErrors {
		return nil
	}
	return err
}
