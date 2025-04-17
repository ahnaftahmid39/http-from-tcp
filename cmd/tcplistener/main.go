package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
)

func main() {
	tcp, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to setup tcp server")
		os.Exit(1)
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Fprintln(os.Stdout, "error establishing connection")
		}
		// fmt.Fprintln(os.Stdout, "Accepted connection from", conn.RemoteAddr())

		go func() {
			req, err := request.RequestFromReader(conn)
			if err != nil {
				fmt.Fprintln(os.Stdout, "Error handling request", err)
			}
			fmt.Fprintln(os.Stdout, "Request line:")
			fmt.Fprintln(os.Stdout, "- Method:", req.RequestLine.Method)
			fmt.Fprintln(os.Stdout, "- Target:", req.RequestLine.RequestTarget)
			fmt.Fprintln(os.Stdout, "- Version:", req.RequestLine.HttpVersion)

			fmt.Fprintln(os.Stdout, "Headers:")
			for key, val := range req.Headers {
				fmt.Fprintf(os.Stdout, "- %s: %s\n", key, val)
			}
			conn.Close()
			// fmt.Fprintln(os.Stdout, "Connection to", conn.RemoteAddr(), "has been closed")
		}()
	}

}
