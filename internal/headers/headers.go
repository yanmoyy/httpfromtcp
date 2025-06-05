package headers

import (
	"bytes"
	"fmt"
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
	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	h.Set(key, string(value))
	return idx + len(crlf), false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
