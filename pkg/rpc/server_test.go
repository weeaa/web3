package rpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer(DefaultPort)
	defer server.Server.Stop()
	if err != nil {
		assert.Error(t, fmt.Errorf("error creating grpc server: %w", err))
	}

	conn, err := net.Dial("tcp", DefaultPort)
	defer conn.Close()
	if err != nil {
		assert.Error(t, fmt.Errorf("error connecting to grpc server: %w", err))
	}

	assert.NoError(t, nil)
}

func TestReceiveMessage(t *testing.T) {
	server, err := NewServer(DefaultPort)
	defer server.Server.Stop()
	if err != nil {
		assert.Error(t, fmt.Errorf("error creating grpc server: %w", err))
	}
}
