package server_test

import (
	"testing"
	"time"

	"github.com/chapsuk/mserv/server"
	"github.com/stretchr/testify/assert"
	"layeh.com/radius"
)

func TestRadiusServer(t *testing.T) {
	var (
		handler = func() radius.HandlerFunc {
			return func(w radius.ResponseWriter, r *radius.Request) {
				w.Write(r.Response(radius.CodeAccessAccept))
			}
		}

		rsrv = &radius.PacketServer{
			Addr:               ":8086",
			Network:            "udp",
			SecretSource:       radius.StaticSecretSource([]byte("secret")),
			Handler:            handler(),
			InsecureSkipVerify: true,
		}
	)

	srv, err := server.NewRadius(rsrv,
		server.RadiusSkipErrors(false),
		server.RadiusShutdownTimeout(time.Second))

	assert.NoError(t, err)
	assert.NotNil(t, srv)

	assert.NoError(t, srv.Start())
	assert.NoError(t, srv.Stop())
}
