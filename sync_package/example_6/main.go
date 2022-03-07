package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Using Sync pool
// Create a pool of bytes.buffer which can be reused.
// rather than creating new buffer instance each time.

// New function will called when there is no instance available
// in the buffer Pool
var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Allocating new bytes.Buffer")
		return new(bytes.Buffer)
	},
}

func log(w io.Writer, debug string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(debug)
	b.WriteString("\n")

	w.Write(b.Bytes())
	// once we used buffer, put back them to the poolcd ..
	// to reuse them in future
	bufPool.Put(b)
}
func main() {
	log(os.Stdout, "first-log")
	log(os.Stdout, "sceond-log")
	log(os.Stdout, "third-log")
}
