package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/yanmoyy/httpfromtcp/internal/request"
	"github.com/yanmoyy/httpfromtcp/internal/response"
	"github.com/yanmoyy/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	_ = w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	_ = w.WriteHeaders(h)
	_, _ = w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	_ = w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	_ = w.WriteHeaders(h)
	_, _ = w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	_ = w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	_ = w.WriteHeaders(h)
	_, _ = w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get("https://httpbin.org/" + target)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	_ = w.WriteStatusLine(response.StatusCodeSuccess)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Override("Transfer-Encoding", "chunked")
	_ = w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Println("Read ", n, " Bytes")
		if n > 0 {
			_, err := w.WriteChunkBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body: ", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body: ", err)
			break
		}
	}

	_, err = w.WriteChunkBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done: ", err)
	}
}
