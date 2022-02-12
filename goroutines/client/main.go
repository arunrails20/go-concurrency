package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// connect to server on localhosr port 3000

	conn, err := net.Dial("tcp", "localhost:3000")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	//Copy the Response from the server to Stdout
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
