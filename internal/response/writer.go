package response

import (
	"fmt"
	"io"

	"github.com/yanmoyy/httpfromtcp/internal/headers"
)

type WriterState int

const (
	writerStateStatusLine WriterState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	state  WriterState
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state:  writerStateStatusLine,
		writer: w,
	}
}
func (w *Writer) setState(state WriterState) {
	w.state = state
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.state)
	}
	defer w.setState(writerStateHeaders)
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return fmt.Errorf("cannot write Headers in state %d", w.state)
	}
	defer w.setState(writerStateBody)
	for k, v := range headers {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w.writer, "\r\n")
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	return w.writer.Write(p)
}
