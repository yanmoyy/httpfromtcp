package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	// if no CRLF: need more data
	if idx == -1 {
		return 0, false, nil
	}
	// if CRLF is on the first idx: end of headers
	if idx == 0 {
		return len(crlf), true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}
	key = strings.TrimSpace(key)
	value := bytes.TrimSpace(parts[1])

	if !validToken([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}
	h.Set(key, string(value))
	return idx + len(crlf), false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if val, ok := h[key]; ok {
		value = strings.Join([]string{val, value}, ", ")
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func validToken(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}
	return slices.Contains(tokenChars, c)
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
