package server

import (
	"io"

	"github.com/yanmoyy/httpfromtcp/internal/request"
	"github.com/yanmoyy/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (hErr *HandlerError) Write(w io.Writer) {
	_ = response.WriteStatusLine(w, hErr.StatusCode)
	messageByte := []byte(hErr.Message)
	headers := response.GetDefaultHeaders(len(messageByte))
	_ = response.WriteHeaders(w, headers)
	_, _ = w.Write(messageByte)
}
