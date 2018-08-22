package mserv

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// HTTPServer wrapper of http.Server
type HTTPServer struct {
	shutdownTimeout time.Duration
	server          *http.Server
}

// NewHTTPServer returns new http.Server wrapper
func NewHTTPServer(shutdownTimeout time.Duration, server *http.Server) (Server, error) {
	if server == nil {
		return nil, errors.New("missing http.Server, skip")
	}

	if server.Addr == "" {
		return nil, errors.New("missing bind address for http.Server, skip")
	}

	return &HTTPServer{
		shutdownTimeout: shutdownTimeout,
		server:          server,
	}, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Start http server in goroutine
// write fatal msg by log if cant start server
func (h *HTTPServer) Start() error {
	if h == nil {
		return errors.New("missing handler for http.Server")
	}

	ln, err := net.Listen("tcp", h.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		if err := h.server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("start http server error: %s", err)
			}
		}
	}()

	return nil
}

// Stop http server with timeout
func (h *HTTPServer) Stop() error {
	if h == nil {
		return nil
	}

	if h.shutdownTimeout == 0 {
		if err := h.server.Close(); err != nil {
			return err
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.shutdownTimeout)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop http server error: %s", err)
	}

	return nil
}
