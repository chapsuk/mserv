package mserv

import "go.uber.org/dig"

// Server interface
type Server interface {
	Start()
	Stop()
}

// MultiServer is servers aggregator
type MultiServer struct {
	servers []Server
}

// New yiield new multiple servers instance pointer
func New(servers ...Server) Server {
	m := &MultiServer{}
	for _, s := range servers {
		if s != nil {
			m.servers = append(m.servers, s)
		}
	}
	return m
}

// Start servers
func (ms *MultiServer) Start() {
	for _, s := range ms.servers {
		if s != nil {
			s.Start()
		}
	}
}

// Stop multiple servers and return concatenated error
func (ms *MultiServer) Stop() {
	for _, s := range ms.servers {
		if s != nil {
			s.Stop()
		}
	}
}

// DigMultiServerParams inject dig providers with `group:"mserv"` tag
type DigMultiServerParams struct {
	dig.In
	Servers []Server `group:"mserv"`
}

// NewByDig returns new multi server instance from incomming params
func NewByDig(p DigMultiServerParams) Server {
	return New(p.Servers...)
}
