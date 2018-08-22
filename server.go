package mserv

import "errors"

// Server interface
type Server interface {
	Start() error
	Stop() error
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
func (ms *MultiServer) Start() error {
	if log == nil {
		return errors.New("missing logger, you must call mserv.SetLogger")
	}

	for _, s := range ms.servers {
		if s == nil {
			continue
		}

		if err := s.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop multiple servers and return concatenated error
func (ms *MultiServer) Stop() error {
	for _, s := range ms.servers {
		if s == nil {
			continue
		}
		if err := s.Stop(); err != nil {
			return err
		}
	}

	return nil
}
