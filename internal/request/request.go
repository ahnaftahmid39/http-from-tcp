package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	initialized = iota
	done
)

type Request struct {
	state       int
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(requestLine string) (RequestLine, error) {
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		return RequestLine{}, errors.New("request Line is not good")
	}

	httpVersionParts := strings.Split(parts[2], "/")
	if len(httpVersionParts) < 2 {
		return RequestLine{}, errors.New("something wrong with http version")
	}
	return RequestLine{
		HttpVersion:   httpVersionParts[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read all the request %v", err)
	}

	lines := strings.Split(string(req), "\r\n")

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil

}
