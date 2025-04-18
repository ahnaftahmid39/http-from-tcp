package server

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
	"github.com/ahnaftahmid39/http-from-tcp/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
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
		if s.closed.Load() {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.closed.Load() {
				fmt.Fprintln(os.Stdout, "error establishing connection")
			}
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	hErr := s.handler(buffer, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	b := buffer.Bytes()
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(len(b)))
	conn.Write(b)

}

func Serve(port int, handler Handler) (*Server, error) {
	tcp, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to setup tcp server")
		os.Exit(1)
	}

	server := &Server{
		listener: tcp,
		closed:   atomic.Bool{},
		handler:  handler,
	}

	go server.listen()

	return server, nil
}
