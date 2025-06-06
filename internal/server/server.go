package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/yanmoyy/httpfromtcp/internal/request"
	"github.com/yanmoyy/httpfromtcp/internal/response"
)

// Contains the state of server
type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := Server{
		handler:  handler,
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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := fmt.Appendf(nil, "Error parsing request: %v", err)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}
