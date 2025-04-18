package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil && !s.closed.Load() {
			fmt.Fprintln(os.Stdout, "error establishing connection")
		}
		fmt.Fprintln(os.Stdout, "Accepted connection from", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error handling request", err)
	}
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")
	fmt.Fprintln(os.Stdout, "Connection to", conn.RemoteAddr(), "has been closed")
}

func Serve(port int) (*Server, error) {
	tcp, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to setup tcp server")
		os.Exit(1)
	}

	server := &Server{
		listener: tcp,
		closed:   atomic.Bool{},
	}

	go server.listen()

	return server, nil
}
