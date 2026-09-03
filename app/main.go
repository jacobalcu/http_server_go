package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

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
	lines := strings.Split(reqString, "\r\n")
	reqLine := lines[0]

	// Split req line by spaces
	reqParts := strings.Split(reqLine, " ")

	// Check for METHOD, TARGET, VERSION
	if len(reqParts) < 3 {
		fmt.Println("Malformed request line")
		return
	}

	// Should return str([]bytes) representing path
	target := reqParts[1]

	response := "HTTP/1.1 404 Not Found\r\n\r\n"

	if target == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response: ", err)
		os.Exit(1)
	}
	
}
