package mserv

import (
	multierror "github.com/hashicorp/go-multierror"
)

// Server interface
type Server interface {
	Start() error
	Stop() error
}

// MultiServer is servers aggregator
type MultiServer struct {
	servers []Server
}

// New returns new multiple servers instance, skip nil servers
func New(servers ...Server) Server {
	m := &MultiServer{}
	for _, s := range servers {
		if s != nil {
			m.servers = append(m.servers, s)
		}
	}
	return m
}

// Start calls Start function for each server in group, returns first error when happen
func (ms *MultiServer) Start() error {
	for _, s := range ms.servers {
		if err := s.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Stop multiple servers and returns multierrr
func (ms *MultiServer) Stop() error {
	var rerr error
	for _, s := range ms.servers {
		if err := s.Stop(); err != nil {
			rerr = multierror.Append(rerr, err)
		}
	}
	return rerr
}
