package headers

import (
	"bytes"
	"fmt"
	"strings"

	"slices"

	"github.com/ahnaftahmid39/http-from-tcp/internal/constants"
)

type Headers map[string]string

func validateHeaderKey(headerKey []byte) error {
	if len(headerKey) < 1 {
		return fmt.Errorf("error: length of header key must atleast 1 found: %v", len(headerKey))
	}

	specialChars := []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	for _, ch := range headerKey {
		// Check if character is alphanumeric or includes special characters
		if ('a' <= ch && ch <= 'z') ||
			('A' <= ch && ch <= 'Z') ||
			('0' <= ch && ch <= '9') ||
			slices.Contains(specialChars, ch) {
			continue
		}
		return fmt.Errorf("error: invalid character '%c' found in header key: %s", ch, string(headerKey))
	}
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(constants.CRLF))

	// CRLF not found yet so cannot parse yet
	if crlfIdx == -1 {
		return 0, false, nil
	}

	// CRLF at start detected to header parsing is complete
	if crlfIdx == 0 {
		return crlfIdx + len(constants.CRLF), true, nil
	}

	colonIdx := bytes.Index(data, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("error parsing header, colon not found in field-line: %s", string(data))
	}

	// header key must not have trailing whitespace
	// colonIdx == 0 will handled in validateHeaderKey func
	if colonIdx > 0 && data[colonIdx-1] == ' ' {
		return 0, false, fmt.Errorf("error parsing header, found whitespace before colon in field-line: %s", string(data))
	}

	headerKey := bytes.Trim(data[:colonIdx], " ")
	headerValue := bytes.Trim(data[colonIdx+1:crlfIdx], " ")

	err = validateHeaderKey(headerKey)
	if err != nil {
		return 0, false, err
	}

	headerKeyStr := strings.ToLower(string(headerKey))
	headerValueStr := string(headerValue)
	h.Set(headerKeyStr, headerValueStr)

	return crlfIdx + len(constants.CRLF), false, nil
}

func (h Headers) Set(key, val string) {
	// If already exists, append it with comma separated
	if existing, ok := h[key]; ok {
		h[key] = strings.Join([]string{existing, val}, ", ")
	} else {
		h[key] = val
	}
}

func (h Headers) Get(key string) string {
	return h[key]
}

func NewHeaders() Headers {
	return make(Headers)
}
