package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

// Contains the state of server
type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{
		listener: l,
	}
	go s.listen()
	return &s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	data := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello World!\n"
	_, _ = conn.Write([]byte(data))
}
