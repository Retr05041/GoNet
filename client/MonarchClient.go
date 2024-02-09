package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Display server response (life cycle)
// Should be text coming in from others in the same channel (global)
func recieveFromChannel(serverReader io.Reader) {
	inbuf := make([]byte, 1024)
	for {
		n, err := serverReader.Read(inbuf[:])
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(string(inbuf[:n]))
	}
}

// Runner for client
// Connects to specified server, Creates a Scanner and Reader, and continuesly scans and reads
func main() {

	// Connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Scanner for user input - Reader for server responses
	serverScanner := bufio.NewScanner(os.Stdin)
	clientReader := bufio.NewReader(conn)
	go recieveFromChannel(clientReader)

	// Scan user input
	for serverScanner.Scan() {
		if serverScanner.Err() != nil {
			log.Println(serverScanner.Err())
		}
		data := serverScanner.Text()

		dataAsBytes := []byte(data + "\n")
		_, err := conn.Write(dataAsBytes)
		if err != nil {
			fmt.Println("Failed to send data to the server!")
			break
		}
	}
}
