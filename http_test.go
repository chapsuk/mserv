package mserv_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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

func TestHTTPWithTLSConfig(t *testing.T) {
	addr := testAddr()
	uri := "https://" + addr

	var (
		caCert               = readFile(t, "certs/ca.pem")
		serverCertPem        = readFile(t, "certs/server.pem")
		serverKeyPem         = readFile(t, "certs/server-key.pem")
		clientCertPem        = readFile(t, "certs/client.pem")
		clientKeyPem         = readFile(t, "certs/client-key.pem")
		clientUnknownCertPem = readFile(t, "certs/client-unkwn.pem")
		clientUnknownKeyPem  = readFile(t, "certs/client-unkwn-key.pem")
	)

	serverCert, err := tls.X509KeyPair(serverCertPem, serverKeyPem)
	require.NoError(t, err)

	clientCert, err := tls.X509KeyPair(clientCertPem, clientKeyPem)
	require.NoError(t, err)

	clientUnknownCert, err := tls.X509KeyPair(clientUnknownCertPem, clientUnknownKeyPem)
	require.NoError(t, err)

	cas := x509.NewCertPool()
	cas.AppendCertsFromPEM(caCert)

	serverTlsCfg := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    cas,
	}

	srv := mserv.NewHTTPServer(
		&http.Server{
			Addr: addr,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			}),
		},
		mserv.HTTPWithTLSConfig(serverTlsCfg),
	)
	require.NoError(t, srv.Start())

	cas = x509.NewCertPool()
	cas.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      cas,
			},
		},
	}

	res, err := client.Get(uri)
	require.NoError(t, err)

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.NoError(t, res.Body.Close())
	assert.NoError(t, srv.Stop())

	clientUnkwnon := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientUnknownCert},
				RootCAs:      cas,
			},
		},
	}

	res, err = clientUnkwnon.Get(uri)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func readFile(t *testing.T, file string) []byte {
	b, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return b
}
