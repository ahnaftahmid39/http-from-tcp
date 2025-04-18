package server

import (
	"io"

	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
	"github.com/ahnaftahmid39/http-from-tcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h HandlerError) Error() string {
	return h.Message
}

func (h HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(h.Message)))
	w.Write([]byte(h.Message))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
