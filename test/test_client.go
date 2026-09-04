package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Sim 3 concurrent clients
	for i := 1; i <= 3; i++ {
		go func(id int) {
			// Connect to server
			conn, err := net.Dial("tcp", "localhost:4221")
			if err != nil {
				fmt.Printf("Client %d failed to connect: %v\n", id, err)
				return
			}
			defer conn.Close()
			fmt.Printf("Client %d connected! Sleeping for 3 seconds...\n", id)

			// Wait 3 seconds before sending req
			time.Sleep(3 * time.Second)

			// Send req
			conn.Write([]byte("GET / HTTP/1.1\r\n"))

			// Read response
			buf := make([]byte, 1024)
			n, _ := conn.Read(buf)
			fmt.Printf("Client %d received: %q\n", id, string(buf[:n]))
		}(i)
	}

	// Keep test alive long enough for goroutine to finish
	time.Sleep(5 * time.Second)
	fmt.Println("Test complete")
}