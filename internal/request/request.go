package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ahnaftahmid39/http-from-tcp/internal/constants"
	"github.com/ahnaftahmid39/http-from-tcp/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	state requestState
	// if contentLength is -1 then it means Content-Length header is not present
	contentLength int
	RequestLine   RequestLine
	Headers       headers.Headers
	Body          []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequest() *Request {
	return &Request{
		state:       requestStateInitialized,
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		Body:        nil,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.state == requestStateInitialized {
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.state = requestStateParsingHeaders
		return n, nil
	}
	if r.state == requestStateParsingHeaders {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return n, err
		}
		if done {
			r.state = requestStateParsingBody
		}
		return n, nil
	}
	if r.state == requestStateParsingBody {
		if r.Body == nil {
			r.Body = make([]byte, 0)
			contentLengthHeader := r.Headers.Get("Content-Length")
			if contentLengthHeader == "" {
				r.contentLength = 0
			} else {
				contentLength, err := strconv.Atoi(contentLengthHeader)
				if err != nil {
					return 0, fmt.Errorf("error: content length header invalid value, %v", err)
				}
				r.contentLength = contentLength
			}
		}
		if r.contentLength == 0 {
			r.state = requestStateDone
			return len(data), nil
		}
		r.Body = append(r.Body, data...)
		if len(r.Body) > r.contentLength {
			return len(data), fmt.Errorf("error: body length is more than specified content-length")
		}
		return len(data), nil
	}
	if r.state == requestStateDone {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state %v", r.state)
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, constants.BUFFER_SIZE)
	readToIndex := 0

	req := NewRequest()

	for req.state != requestStateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		readToIndex += n
		if err != nil && errors.Is(err, io.EOF) {
			if req.contentLength != len(req.Body) {
				return nil, fmt.Errorf("error: content-length != len(req.Body), %v != %v respectively", req.contentLength, len(req.Body))
			}
			req.state = requestStateDone
			break
		}
		if err != nil {
			return nil, err
		}
		np, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[np:])
		readToIndex -= np
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(constants.CRLF))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + len(constants.CRLF), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
