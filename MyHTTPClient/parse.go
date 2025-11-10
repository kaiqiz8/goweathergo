package myhttpclient

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	Proto      string
}

func getHeader(headers map[string][]string, name string) ([]string, bool) {
	lname := strings.ToLower(name)
	for k, v := range headers {
		if strings.ToLower(k) == lname {
			return v, true
		}
	}
	return nil, false
}

// parseChunkedBody reads a chunked Transfer-Encoding body
func parseChunkedBody(reader *bufio.Reader) ([]byte, error) {
	var body []byte
	for {
		// Read the chunk size line
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk size: %w", err)
		}
		line = strings.TrimSpace(line)
		size, err := strconv.ParseInt(line, 16, 64) // Hexadecimal size
		if err != nil {
			return nil, fmt.Errorf("invalid chunk size: %w", err)
		}
		if size == 0 {
			// Read trailer headers (if any) and final CRLF
			_, err = reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("failed to read final CRLF: %w", err)
			}
			break
		}

		// Read the chunk data
		chunk := make([]byte, size)
		_, err = reader.Read(chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk: %w", err)
		}
		body = append(body, chunk...)

		// Read the trailing CRLF
		_, err = reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk trailing CRLF: %w", err)
		}
	}
	return body, nil
}

func ParseResponse(reader *bufio.Reader) (Response, error) {
	// Parse the status line
	statusLine, err := reader.ReadString('\n') // Example: HTTP/1.1 200 OK\r\n
	if err != nil {
		return Response{}, fmt.Errorf("failed to read status line: %w", err)
	}
	statusParts := strings.SplitN(strings.TrimSpace(statusLine), " ", 3)
	if len(statusParts) < 2 {
		return Response{}, fmt.Errorf("malformed status line: %s", statusLine)
	}
	proto := statusParts[0]
	if !strings.HasPrefix(proto, "HTTP/") {
		return Response{}, fmt.Errorf("invalid protocol: %s", proto)
	}

	statusCode, err := strconv.Atoi(statusParts[1]) // strconv.Atoi converts "200" (string) → 200 (int)
	if err != nil {
		return Response{}, fmt.Errorf("invalid status code: %w", err)
	}

	// Read headers
	headers := make(map[string][]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return Response{}, fmt.Errorf("failed to read headers: %w", err)
		}
		if line == "\r\n" {
			break // End of headers
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			key := strings.ToLower(strings.TrimSpace(headerParts[0])) // case-insensitive header names
			value := strings.TrimSpace(headerParts[1])
			headers[key] = append(headers[key], value)
		}
	}

	// Determine body length and read body
	var body []byte

	// Check chunked transfer
	if encodings, ok := getHeader(headers, "Transfer-Encoding"); ok {
		for _, encoding := range encodings {
			if strings.ToLower(encoding) == "chunked" {
				body, err = parseChunkedBody(reader)
				if err != nil {
					return Response{}, fmt.Errorf("failed to parse chunked body: %w", err)
				}
			}
		}
	}

	if lengths, ok := getHeader(headers, "Content-Length"); ok && len(lengths) > 0 {
		_, err := strconv.Atoi(lengths[0])
		if err != nil {
			return Response{}, fmt.Errorf("invalid Content-Length: %w", err)
		}
		body = make([]byte, 0)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return Response{}, fmt.Errorf("failed to read body: %w", err)
		}
	}

	if len(body) == 0 && statusCode != 204 && statusCode != 304 && statusCode/100 != 1 {
		body, err = reader.ReadBytes(0) // Read until EOF
		if err != nil && err.Error() != "EOF" {
			return Response{}, fmt.Errorf("failed to read body: %w", err)
		}
	}

	return Response{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
		Proto:      proto,
	}, nil
}
