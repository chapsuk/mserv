package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/chapsuk/grace"
	"github.com/chapsuk/mserv"
	"github.com/chapsuk/mserv/examples/helloworld"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serverFunc = func() (mserv.Server, error)

// to simplify error catch
var servers = []serverFunc{
	// pprof
	func() (mserv.Server, error) {
		return mserv.NewHTTPServer(time.Second, &http.Server{
			Addr:    ":8081",
			Handler: http.DefaultServeMux,
		})
	},

	// prometheus
	func() (mserv.Server, error) {
		return mserv.NewHTTPServer(time.Second, &http.Server{
			Addr:    ":8082",
			Handler: promhttp.Handler(),
		})
	},

	// gin
	func() (mserv.Server, error) {
		return mserv.NewHTTPServer(5*time.Second, &http.Server{
			Addr:         ":8083",
			Handler:      ginApp(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		})
	},

	// echo
	func() (mserv.Server, error) {
		return mserv.NewHTTPServer(5*time.Second, &http.Server{
			Addr:         ":8084",
			Handler:      echoApp(),
			TLSConfig:    &tls.Config{ /**todo**/ },
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	},

	// grpc
	func() (mserv.Server, error) {
		return mserv.NewGRPCServer(":8085", grpcServer())
	},
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	mserv.SetLogger(logger)

	serves := make([]mserv.Server, 0, len(servers))

	for _, serve := range servers {
		s, err := serve()
		if err != nil {
			panic(err)
		}

		serves = append(serves, s)
	}

	s := mserv.New(serves...)

	logger.Println("Start servers")
	if err := s.Start(); err != nil {
		panic(err)
	}

	<-grace.ShutdownContext(context.Background()).Done()
	logger.Printf("Received stop signal")

	logger.Println("Stop servers")
	if err := s.Stop(); err != nil {
		panic(err)
	}
	logger.Println("Gracefully stopped")
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
	gin.SetMode(gin.ReleaseMode)

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
