package main 

import (
    "fmt"
    "log"
    "net"
	"strings"
)

func HandleConnection(conn net.Conn) {
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

		if strings.TrimRight(data, "\n") == ":quit" {
			data := ":quit"
			_, err = conn.Write([]byte(data))
			conn.Close()
			return
		}

        _, err = conn.Write([]byte(data))
        if err != nil {
            log.Println(err)
            return
        }
    }
}


func main() {
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
		go HandleConnection(conn)
	}
}
















