# Mserv

[![Build Status](https://travis-ci.org/chapsuk/mserv.svg?branch=master)](https://travis-ci.org/chapsuk/mserv)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchapsuk%2Fmserv.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchapsuk%2Fmserv?ref=badge_shield)

Package for grouping many servers (http, grpc, ...) and start|stop by single command.
Aims to simplify startup process and controll group of servers.

## Example

*Start 3 http servers: pprof, prometheus, api.*


```go
// set  skipErros option for not critical component
srv1 := mserv.NewHTTPServer(&http.Server{Addr: "8080", http.DefaultServeMux}, mserv.HTTPSkipErrors(true))
srv2 := mserv.NewHTTPServer(&http.Server{Addr: "8081", promhttp.Handler()})
// set shutdown timeout
srv3 := mserv.NewHTTPServer(&http.Server{Addr: "8082", api()}, mserv.HTTPShutdownTimeout(5*time.Second))

// the servers order is keeped at startup
srvs := mserv.New(srv1, srv2, srv3)
err  := srvs.Start()
// ... do work ...
err := srvs.Stop()
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchapsuk%2Fmserv.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchapsuk%2Fmserv?ref=badge_large)