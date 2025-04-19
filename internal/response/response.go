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

func writeStatusLine(w io.Writer, statusCode StatusCode) error {
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
		"content-length": strconv.Itoa(contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
	}
}

func writeHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		fmt.Fprintf(w, "%s: %s\r\n", key, val)
	}
	fmt.Fprintf(w, "\r\n")
	return nil
}

type WriterState int

const (
	WriterStateStatusLine WriterState = iota
	WriterStateHeaders
	WriterStateBody
)

type Writer struct {
	Writer io.Writer

	state WriterState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateStatusLine {
		return fmt.Errorf("error: invalid state to write status line")
	}
	writeStatusLine(w.Writer, statusCode)
	w.state = WriterStateHeaders
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStateHeaders {
		return fmt.Errorf(("error: invalid state to write headers"))
	}

	finalHeaders := GetDefaultHeaders(0)
	for key, val := range headers {
		finalHeaders.Override(key, val)
	}

	fmt.Println(finalHeaders)

	writeHeaders(w.Writer, finalHeaders)
	w.state = WriterStateBody

	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("error: invalid state to write body")
	}

	return w.Writer.Write(p)
}
