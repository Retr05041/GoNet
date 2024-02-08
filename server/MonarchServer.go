package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type server struct {
	clients []io.Writer
}

func (s *server) addClient(c net.Conn) {
	s.clients = append(s.clients, c)
}

func (s *server) writeAll(data string) {
	for _, cl := range s.clients {
		_, err := cl.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		userInput, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		data := string(buf[:userInput])
		fmt.Printf("Received: %s", data)

		s.writeAll(data)
	}
}

func main() {
	fmt.Println("Starting server...")
	srv := &server{}

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port 8000...")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		srv.addClient(conn)

		go srv.HandleConnection(conn)
	}
}
