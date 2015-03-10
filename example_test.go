package tcptest_test

import (
	"fmt"
	"github.com/stvp/tcptest"
	"log"
	"net"
	"time"
)

func ExampleServer() {
	server, err := tcptest.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		conn, err := net.Dial("tcp", server.Address())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(conn, "hello\nworld!\n")
		conn.Close()
	}()

	server.WaitForLines(2, time.Second)

	if server.ReceivedLine("hello") && server.Received("world") {
		fmt.Printf("received!")
	}
	// Output: received!

	server.Close()
}
