package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

// server: Struct for holding server info
type server struct {
	// List of clients
	clients []client
}

// client: Struct for holding client info
type client struct {
	username         string    // Clients given username for global communication
	clientWriter     io.Writer // Uses conn as it's io.Writer
	clientConnection net.Conn  // Connection interface
	clientUUID       uuid.UUID // Unique UUID for each client
}

// addClient: Add client to list of clients on server
func (s *server) AddClient(c client) int {
	s.clients = append(s.clients, c)
  return 1
}

// writeAll: Send a string to every client the server knows excluding the providingClient
func (s *server) WriteAll(providingClient uuid.UUID, data string) {
	for _, cl := range s.clients {
		if cl.clientUUID != providingClient {
			_, err := cl.clientWriter.Write([]byte(data))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// HandleConnections: Main life cycle for every client connection
func (s *server) HandleConnection(c client) {
	defer c.clientConnection.Close()
	// defer fmt.Println("Connection closed with client.")

	for {
		clientData, err := bufio.NewReader(c.clientConnection).ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		cleanedData := strings.TrimSpace(strings.TrimRight(string(clientData), "\n"))
		fmt.Println(cleanedData)
		// fmt.Println(cleanedData)

		go s.WriteAll(c.clientUUID, cleanedData)
	}
}

// main: Runner for server
func main() {
	fmt.Println("Starting server...")
	srv := &server{}

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port 8000...")
	defer listener.Close()

	// Accept incoming connections forever
	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Create a client
		newClient := client{
			clientWriter:     conn,
			clientConnection: conn,
			clientUUID:       uuid.New(),
		}

		// fmt.Println("Connection made with client.")
		// Add client to the client list then begin client life cycle
		srv.AddClient(newClient)
		go srv.HandleConnection(newClient)
	}
}
