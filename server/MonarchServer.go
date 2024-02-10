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

// Struct for server based info
type server struct {
	clients []client
}

// Struct for individual Client info
type client struct {
	username         string
	clientWriter     io.Writer
	clientConnection net.Conn
	clientUUID       uuid.UUID
}

// Add a client to the server struct
func (s *server) addClient(c client) {
	s.clients = append(s.clients, c)
}

// Send data to every client in the channel
func (s *server) writeAll(providingClient uuid.UUID, data string) {
	for _, cl := range s.clients {
		if cl.clientUUID != providingClient {
			_, err := cl.clientWriter.Write([]byte(data))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// Handle incoming connections
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

		go s.writeAll(c.clientUUID, cleanedData)
	}
}

// Runner for the server
func main() {
	fmt.Println("Starting server...")
	srv := &server{}

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port 8000...")
	defer listener.Close()

	// Do this forever
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
		srv.addClient(newClient)
		go srv.HandleConnection(newClient)
	}
}
