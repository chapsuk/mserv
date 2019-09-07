package mserv_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/chapsuk/mserv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var startPort = 11101
var testAddr = func() string {
	startPort++
	return fmt.Sprintf("127.0.0.1:%d", startPort)
}

func TestMServerNilServer(t *testing.T) {
	srv := mserv.New(nil, nil, nil)
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestMServerStartStop(t *testing.T) {
	var (
		grpcAddr = testAddr()
		httpAddr = testAddr()
	)

	srv := mserv.New(
		mserv.NewGRPCServer(grpcAddr, grpc.NewServer()),
		mserv.NewHTTPServer(&http.Server{Addr: httpAddr}),
	)

	assert.NoError(t, srv.Start())
	assert.True(t, hasListener(grpcAddr))
	assert.True(t, hasListener(httpAddr))
	assert.NoError(t, srv.Stop())
}

func TestMServerFailStart(t *testing.T) {
	var httpAddr = testAddr()

	srv := mserv.New(
		mserv.NewGRPCServer("foo", grpc.NewServer()),
		mserv.NewHTTPServer(&http.Server{Addr: httpAddr}),
	)

	assert.Error(t, srv.Start())
	assert.False(t, hasListener(httpAddr))
	assert.NoError(t, srv.Stop())
}

func TestMserverStopErr(t *testing.T) {
	srv := mserv.New(
		&fakeServer{stopErr: errors.New("test")},
	)

	assert.NoError(t, srv.Start())
	require.Error(t, srv.Stop())
}

func hasListener(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	return err == nil
}

type fakeServer struct {
	startErr error
	stopErr  error
}

func (s *fakeServer) Start() error {
	return s.startErr
}

func (s *fakeServer) Stop() error {
	return s.stopErr
}
