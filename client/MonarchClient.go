package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type clientInfo struct {
	username string
}

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

	client := &clientInfo{}

	fmt.Printf("Username: ")
	fmt.Scanln(&client.username)

	// Connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(client.username + "\n"))
	if err != nil {
		fmt.Println("Failed to send username to the server!")
	}

	// Scanner for user input - Reader for server responses
	serverReader := bufio.NewReader(conn)
	go recieveFromChannel(serverReader)
	clientScanner := bufio.NewScanner(os.Stdin)

	// Scan user input
	for clientScanner.Scan() {
		if clientScanner.Err() != nil {
			log.Println(clientScanner.Err())
		}
		data := clientScanner.Text()

		dataAsBytes := []byte(data + "\n")
		_, err := conn.Write(dataAsBytes)
		if err != nil {
			fmt.Println("Failed to send data to the server!")
			break
		}
	}
}
