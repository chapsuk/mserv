package mserv_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/chapsuk/mserv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPBadAddr(t *testing.T) {
	srv := mserv.NewHTTPServer(&http.Server{Addr: ""})
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())

	srv = mserv.NewHTTPServer(&http.Server{Addr: "foo"})
	assert.Error(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestHTTPSkipErrors(t *testing.T) {
	srv := mserv.NewHTTPServer(&http.Server{Addr: ""}, mserv.HTTPSkipErrors(true))
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestHTTPStartStop(t *testing.T) {
	httpAddr := testAddr()
	srv := mserv.NewHTTPServer(&http.Server{Addr: httpAddr})
	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}

func TestHTTPNoTimeout(t *testing.T) {
	httpAddr := testAddr()
	srv := mserv.NewHTTPServer(&http.Server{Addr: httpAddr}, mserv.HTTPShutdownTimeout(0))
	assert.NoError(t, srv.Start())
	conn, err := net.DialTimeout("tcp", httpAddr, 20*time.Millisecond)
	require.NoError(t, err)
	assert.NoError(t, srv.Stop())
	assert.NoError(t, conn.Close())
}
