package server_test

import (
	"fmt"
	"testing"

	"github.com/chapsuk/mserv/server"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var startPort = 12101
var testAddr = func() string {
	startPort++
	return fmt.Sprintf("127.0.0.1:%d", startPort)
}

func TestGRPCServerBadAddr(t *testing.T) {
	srv := server.NewGRPC("", grpc.NewServer())
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())

	srv = server.NewGRPC("foo", grpc.NewServer())
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestGRPCServerStartStop(t *testing.T) {
	grpcAddr := testAddr()
	srv := server.NewGRPC(grpcAddr, grpc.NewServer())
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestGRPCServerSkipErrors(t *testing.T) {
	srv := server.NewGRPC(testAddr(), grpc.NewServer(), server.GRPCSkipErrors(true))
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}
