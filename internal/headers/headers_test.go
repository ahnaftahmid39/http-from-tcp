package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host:      localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069", headers["host"])

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers with one upper and one lowercase
	headers = NewHeaders()
	data = []byte("content-length:0\r\nContent-Length:42069\r\n\r\n")
	n1, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 18, n1)
	assert.False(t, done)
	n2, done, err := headers.Parse(data[n1:])
	require.NoError(t, err)
	assert.Equal(t, n2, len(data)-2-n1) // total - last crlf and first header
	assert.False(t, done)
	assert.Equal(t, "0, 42069", headers["content-length"])
	assert.Equal(t, 1, len(headers))

	// Test: Valid done (continuing from previous)
	n3, done, err := headers.Parse(data[n1+n2:])
	require.NoError(t, err)
	assert.Equal(t, n3, 2)
	assert.True(t, done)

	// Test: Invalid header key contains invalid character
	headers = NewHeaders()
	data = []byte("Hos@t: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Contains(t, err.Error(), "invalid character")

	// Test: Invalid header key length = 0
	headers = NewHeaders()
	data = []byte(": localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Contains(t, err.Error(), "length of header key")

	// Test: Valid header key with some special characters
	headers = NewHeaders()
	data = []byte("X-Custom-Header-~!#$%&'*+-.^_`|~: value123\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)
	assert.Equal(t, "value123", headers["x-custom-header-~!#$%&'*+-.^_`|~"])

}
