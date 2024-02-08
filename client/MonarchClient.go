package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func recvLoop(r io.Reader) {
	inbuf := make([]byte, 1024)
	for {
		n, err := r.Read(inbuf[:])
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(string(inbuf[:n]))
	}
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(os.Stdin)
	rd := bufio.NewReader(conn)
	go recvLoop(rd)

	for sc.Scan() {
		if sc.Err() != nil {
			log.Println(sc.Err())
		}
		txt := sc.Text()

		b := []byte(txt + "\n")
		_, err := conn.Write(b)
		if err != nil {
			fmt.Println("Failed to send data to the server!")
			break
		}
	}
}
