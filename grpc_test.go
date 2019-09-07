package mserv_test

import (
	"testing"

	"github.com/chapsuk/mserv"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGRPCServerBadAddr(t *testing.T) {
	srv := mserv.NewGRPCServer("", grpc.NewServer())
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())

	srv = mserv.NewGRPCServer("foo", grpc.NewServer())
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestGRPCServerStartStop(t *testing.T) {
	grpcAddr := testAddr()
	srv := mserv.NewGRPCServer(grpcAddr, grpc.NewServer())
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestGRPCServerSkipErrors(t *testing.T) {
	srv := mserv.NewGRPCServer(testAddr(), grpc.NewServer(), mserv.GRPCSkipErrors(true))
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}
