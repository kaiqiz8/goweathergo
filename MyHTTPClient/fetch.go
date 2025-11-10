package myhttpclient

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
}

// type Response struct {
// 	StatusCode int
// 	Headers    map[string]string
// 	Body       []byte
// 	Proto      string
// }

func Fetch(req Request) (Response, error) {
	// Parse the URL
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return Response{StatusCode: 404}, fmt.Errorf("invalid URL: %w", err)
	}
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "80" // Default to port 80 for HTTP
	}
	address := net.JoinHostPort(host, port)

	// Open a TCP connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return Response{}, fmt.Errorf("failed to connect to %s server: %w", address, err)
	}
	defer conn.Close() // ensures the connection is closed when this function exits, this function refers to Fetch

	// Compose the HTTP request
	path := parsedURL.RequestURI() // returns the path and query string of the URL
	// example: URL = "http://example.com/path?query=1" -> path (or RequestURI()) = "/path?query=1"
	if path == "" {
		path = "/"
	}
	requestLine := fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, path)
	headers := fmt.Sprintf("Host: %s\r\n", host)
	for key, value := range req.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headers += "Connection: close\r\n\r\n"

	// Send the request
	_, err = conn.Write([]byte(requestLine + headers)) // Converts the request string to bytes and writes it directly to the TCP socket.
	if err != nil {
		return Response{}, fmt.Errorf("failed to send request: %w", err)
	}

	// Read the response
	reader := bufio.NewReader(conn) // Wraps the TCP connection (conn) into a buffered reader
	// Makes it easier to read text line by line (ReadString('\n')).
	// Under the hood, this reads from the network socket.

	return ParseResponse(reader)
}
