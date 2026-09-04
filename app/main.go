package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type RequestLine struct {
	Method string
	Target string
	Version string
}

func parseReqLine(reqLine string) (RequestLine, error) {
	// Split req line by spaces
	reqParts := strings.Fields(reqLine)

	// Check for METHOD, TARGET, VERSION
	if len(reqParts) < 3 {
		return RequestLine{}, errors.New("malformed request line")
	}

	// Instantiate and return in one step
	return RequestLine{
		Method:  reqParts[0],
		Target:  reqParts[1],
		Version: reqParts[2],
	}, nil
}

type HTTPRequest struct {
	RequestLine RequestLine


}



func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Buffer to hold incoming data
	reqBuffer := make([]byte, 1024)
	n, err := conn.Read(reqBuffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	// Convert bytes we read into string
	reqString := string(reqBuffer[:n])

	// Split by lines
	// Will split req line, headers, body
	lines := strings.Split(reqString, "\r\n")
	reqLine := lines[0]

	parsedReq, err := parseReqLine(reqLine)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	
	response := "HTTP/1.1 404 Not Found\r\n\r\n"

	switch {
	case parsedReq.Target == "/":
		response = "HTTP/1.1 200 OK\r\n\r\n"
	case parsedReq.Target == "/user-agent":
		responseHeader := "HTTP/1.1 200 OK"
		contentType := "Content-Type: text/plain"
		var body string
		
		for _, line := range lines {
			// Guard against empty lines, like between headers and body
			if line == "" {
				continue
			}
			// Split on first colon into only two strings
			parts := strings.SplitN(line, ":", 2)

			// Ensure key-val pair
			if len(parts) == 2 {
				headerName := strings.ToLower(strings.TrimSpace(parts[0]))
				// Find User-Agent field
				// Header names are case-insensitive, force lower
				if headerName == "user-agent:" {
					// Grab entire second part, can include spaces
					body = strings.TrimSpace(parts[1])
					break
				}
			}
			
		}
		contentLen := fmt.Sprintf("Content-Length: %d", len(body))
		response = fmt.Sprintf("%s\r\n%s\r\n%s\r\n\r\n%s", responseHeader, contentType, contentLen, body)
	case strings.HasPrefix(parsedReq.Target, "/echo/"):
		responseHeader := "HTTP/1.1 200 OK"
		contentType := "Content-Type: text/plain"
		body := parsedReq.Target[6:]
		contentLength := fmt.Sprintf("Content-Length: %d", len(body))
		response = fmt.Sprintf("%s\r\n%s\r\n%s\r\n\r\n%s", responseHeader, contentType, contentLength, body)
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err)
		os.Exit(1)
	}
	
}
