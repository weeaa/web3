package rpc

import (
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Error creating server: %v", err)
	}

	defer server.Server.Stop()

	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}

	conn.Close()
}

func TestReceiveMessage(t *testing.T) {

}
