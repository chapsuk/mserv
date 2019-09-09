package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"layeh.com/radius"
)

// RadiusOption allows redefine default settings
type RadiusOption func(*RadiusServer)

// RadiusShutdownTimeout sets server shutdown timout, timeout should be great then zero
func RadiusShutdownTimeout(timeout time.Duration) RadiusOption {
	return func(s *RadiusServer) {
		if timeout > 0 {
			s.shutdownTimeout = timeout
		}
	}
}

// RadiusSkipErrors sets skip errors flag
func RadiusSkipErrors(skip bool) RadiusOption {
	return func(s *RadiusServer) {
		s.skipErrors = skip
	}
}

// RadiusServer wraps github.com/layeh/radius packet server
type RadiusServer struct {
	srv             *radius.PacketServer
	shutdownTimeout time.Duration
	skipErrors      bool
}

// NewRadius return server
func NewRadius(s *radius.PacketServer, opts ...RadiusOption) (*RadiusServer, error) {
	// source: https://github.com/layeh/radius/blob/master/server-packet.go#L190-L195
	if s.Handler == nil {
		return nil, errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return nil, errors.New("radius: nil SecretSource")
	}

	srv := &RadiusServer{
		srv:             s,
		skipErrors:      false,
		shutdownTimeout: 3 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv, nil
}

// Start radius server
func (s *RadiusServer) Start() error {
	// source: https://github.com/layeh/radius/blob/master/server-packet.go#L197-L205
	addrStr := ":1812"
	if s.srv.Addr != "" {
		addrStr = s.srv.Addr
	}
	network := "udp"
	if s.srv.Network != "" {
		network = s.srv.Network
	}
	ls, err := net.ListenPacket(network, addrStr)
	if err != nil {
		return s.returnErr(err)
	}

	go func() {
		defer ls.Close()
		if err := s.srv.Serve(ls); err != radius.ErrServerShutdown {
			if !s.skipErrors {
				panic(fmt.Sprintf("radius: unexpected serve error: %s", err))
			}
		}
	}()

	return nil
}

// Stop radius server
func (s *RadiusServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.returnErr(s.srv.Shutdown(ctx))
}

func (s *RadiusServer) returnErr(err error) error {
	if s.skipErrors {
		return nil
	}
	return err
}
