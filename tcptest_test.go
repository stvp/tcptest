package tcptest

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("tcp", server.Address())
	if err != nil {
		t.Fatal(err)
	}

	if server.Received("foo") {
		t.Error("haven't sent anything yet")
	}

	fmt.Fprintf(conn, "cool\nneat\nincomplete")

	err = server.WaitForLines(2, time.Second)
	if err != nil {
		t.Error(err)
	}

	if !server.Received("oo") {
		t.Error("didn't receive 'oo'")
	}
	if !server.ReceivedLine("cool") {
		t.Error("didn't receive line 'cool'")
	}

	err = server.WaitForLines(3, time.Millisecond)
	if err == nil {
		t.Error("expected error, got none")
	}

	if server.Received("incomplete") {
		t.Error("server should only count complete lines as received")
	}
	if server.ReceivedLine("incomplete") {
		t.Error("server should only count complete lines as received")
	}

	expectLines := []string{"cool", "neat"}
	got := server.Lines()
	if !reflect.DeepEqual(got, expectLines) {
		t.Errorf("got lines: %#v", got)
	}

	conn.Close()
	server.Close()
}
