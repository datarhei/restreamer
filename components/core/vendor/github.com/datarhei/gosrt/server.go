package srt

import (
	"errors"
)

// Server is a framework for a SRT server
type Server struct {
	// The address the SRT server should listen on, e.g. ":6001".
	Addr string

	// Config is the configuration for a SRT listener.
	Config *Config

	// HandleConnect will be called for each incoming connection. This
	// allows you to implement your own interpretation of the streamid
	// and authorization. If this is nil, all connections will be
	// rejected.
	HandleConnect AcceptFunc

	// HandlePublish will be called for a publishing connection.
	HandlePublish func(conn Conn)

	// HandlePublish will be called for a subscribing connection.
	HandleSubscribe func(conn Conn)

	ln Listener
}

// ErrServerClosed is returned when the server is about to shutdown.
var ErrServerClosed = errors.New("srt: server closed")

// ListenAndServe starts the SRT server. It blocks until an error happens.
// If the error is ErrServerClosed the server has shutdown normally.
func (s *Server) ListenAndServe() error {
	err := s.Listen()
	if err != nil {
		return err
	}

	return s.Serve()
}

// Listen opens the server listener.
// It returns immediately after the listener is ready.
func (s *Server) Listen() error {
	// Set some defaults if required.
	if s.HandlePublish == nil {
		s.HandlePublish = s.defaultHandler
	}

	if s.HandleSubscribe == nil {
		s.HandleSubscribe = s.defaultHandler
	}

	if s.Config == nil {
		config := DefaultConfig()
		s.Config = &config
	}

	// Start listening for incoming connections.
	ln, err := Listen("srt", s.Addr, *s.Config)
	if err != nil {
		return err
	}

	s.ln = ln

	return err
}

// Serve starts accepting connections. It must be called after Listen().
// It blocks until an error happens.
// If the error is ErrServerClosed the server has shutdown normally.
func (s *Server) Serve() error {
	for {
		// Wait for connections.
		conn, mode, err := s.ln.Accept(s.HandleConnect)
		if err != nil {
			if err == ErrListenerClosed {
				return ErrServerClosed
			}

			return err
		}

		if conn == nil {
			// rejected connection, ignore
			continue
		}

		if mode == PUBLISH {
			go s.HandlePublish(conn)
		} else {
			go s.HandleSubscribe(conn)
		}
	}
}

// Shutdown will shutdown the server. ListenAndServe will return a ErrServerClosed
func (s *Server) Shutdown() {
	if s.ln == nil {
		return
	}

	// Close the listener
	s.ln.Close()
}

func (s *Server) defaultHandler(conn Conn) {
	// Close the incoming connection
	conn.Close()
}
