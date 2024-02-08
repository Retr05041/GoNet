package main 

import (
	"fmt"
	"log"
	"bufio"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Listening on port 8000...")
	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
        message, err :=  bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        fmt.Print("Message Received:", string(message))
        newmessage := strings.ToUpper(message)
        conn.Write([]byte(newmessage + "\n"))
    }
}
















