package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ahnaftahmid39/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	httpVersion := "HTTP/1.1"
	switch statusCode {
	case StatusOk:
		fmt.Fprintf(w, "%s %v %s\r\n", httpVersion, statusCode, "OK")
	case StatusBadRequest:
		fmt.Fprintf(w, "%s %v %s\r\n", httpVersion, statusCode, "Bad Request")
	case StatusInternalServerError:
		fmt.Fprintf(w, "%s %v %s\r\n", httpVersion, statusCode, "Internal Server Error")
	default:
		fmt.Fprintf(w, "%s %v %s\r\n", httpVersion, StatusBadRequest, "")
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		fmt.Fprintf(w, "%s: %s\r\n", key, val)
	}
	fmt.Fprintf(w, "\r\n")
	return nil
}
