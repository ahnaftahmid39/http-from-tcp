package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahnaftahmid39/http-from-tcp/internal/headers"
	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
	"github.com/ahnaftahmid39/http-from-tcp/internal/response"
	"github.com/ahnaftahmid39/http-from-tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			badRequest := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
			w.WriteStatusLine(400)
			w.WriteHeaders(headers.Headers{
				"Content-Type":   "text/html",
				"Content-Length": fmt.Sprintf("%d", len(badRequest)),
			})
			w.WriteBody([]byte(badRequest))

		case "/myproblem":
			serverError := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
			w.WriteStatusLine(500)
			w.WriteHeaders(headers.Headers{
				"Content-Type":   "text/html",
				"Content-Length": fmt.Sprintf("%d", len(serverError)),
			})
			w.WriteBody([]byte(serverError))
		}

		success := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
		w.WriteStatusLine(200)
		w.WriteHeaders(headers.Headers{
			"Content-Type":   "text/html",
			"Content-Length": fmt.Sprintf("%d", len(success)),
		})
		w.WriteBody([]byte(success))
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
