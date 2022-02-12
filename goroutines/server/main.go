package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// program to handle concurrent client connection
	// server listen on protocol TCP and localhost:3000
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	// infinite loop to accept the client connections
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// write response to the connection
		_, err := io.WriteString(c, "response from server\n")
		if err != nil {
			return
		}
	}
	// at interval of one second we are writing the response
	time.Sleep(time.Second)
}
