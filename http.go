package mserv

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// HTTPServerOption allows configure http server optional settings
type HTTPServerOption func(*HTTPServer)

// HTTPShutdownTimeout sets shutdown timeout
func HTTPShutdownTimeout(timeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.shutdownTimeout = timeout
	}
}

// HTTPSkipErrors sets skip errors flag
func HTTPSkipErrors(skip bool) HTTPServerOption {
	return func(s *HTTPServer) {
		s.skipErrors = skip
	}
}

// HTTPWithTLSConfig define tls config for server
func HTTPWithTLSConfig(cfg *tls.Config) HTTPServerOption {
	return func(s *HTTPServer) {
		s.tlsConfig = cfg
	}
}

// HTTPServer wrapper of http.Server
type HTTPServer struct {
	skipErrors      bool
	shutdownTimeout time.Duration
	server          *http.Server
	tlsConfig       *tls.Config
}

// NewHTTPServer returns new http.Server wrapper
func NewHTTPServer(s *http.Server, opts ...HTTPServerOption) Server {
	srv := &HTTPServer{
		server:          s,
		skipErrors:      false,
		shutdownTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

// Start http server
func (h *HTTPServer) Start() error {
	if len(h.server.Addr) == 0 {
		return h.returnErr(errors.New("missing bind addr"))
	}

	ls, err := net.Listen("tcp", h.server.Addr)
	if err != nil {
		return h.returnErr(err)
	}

	if h.tlsConfig != nil {
		ls = tls.NewListener(ls, h.tlsConfig)
	}

	go func() {
		if serr := h.server.Serve(ls); serr != nil && serr != http.ErrServerClosed {
			if !h.skipErrors {
				panic(fmt.Sprintf("http.Serve return unexpected error: %s", serr))
			}
		}
	}()

	return nil
}

// Stop stops http server with timeout
func (h *HTTPServer) Stop() error {
	if h.shutdownTimeout == 0 {
		return h.returnErr(h.server.Close())
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.shutdownTimeout)
	defer cancel()

	return h.returnErr(h.server.Shutdown(ctx))
}

func (h *HTTPServer) returnErr(err error) error {
	if h.skipErrors {
		return nil
	}
	return err
}
