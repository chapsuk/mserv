package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/chapsuk/grace"
	"github.com/chapsuk/mserv"
	"github.com/chapsuk/mserv/examples/helloworld"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "net/http/pprof"
)

func main() {
	s := mserv.New(
		// pprof
		mserv.NewHTTPServer(time.Second, &http.Server{
			Addr:    ":8081",
			Handler: http.DefaultServeMux,
		}),
		// prometheus
		mserv.NewHTTPServer(time.Second, &http.Server{
			Addr:    ":8082",
			Handler: promhttp.Handler(),
		}),
		// gin
		mserv.NewHTTPServer(5*time.Second, &http.Server{
			Addr:         ":8083",
			Handler:      ginApp(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}),
		// echo
		mserv.NewHTTPServer(5*time.Second, &http.Server{
			Addr:         ":8084",
			Handler:      echoApp(),
			TLSConfig:    &tls.Config{ /**todo**/ },
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}),
		// grpc
		mserv.NewGRPCServer(":8085", grpcServer()),
	)

	s.Start()
	<-grace.ShutdownContext(context.Background()).Done()
	s.Stop()
}

func echoApp() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e
}

func ginApp() *gin.Engine {
	router := gin.Default()
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	return router
}

func grpcServer() *grpc.Server {
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	return s
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}
