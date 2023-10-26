package rpc

import (
	"fmt"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
)

const DefaultPort = ":9000"

type ProtoServer struct {
	Server *grpc.Server
}

type ProtoClient struct {
	Conn *grpc.ClientConn
}

func NewServer(port string) (*ProtoServer, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()

	go func() {
		if err = server.Serve(lis); err != nil {
			logger.LogFatal("grpc", err.Error())
		}
	}()

	return &ProtoServer{Server: server}, nil
}

func NewClient() (*ProtoClient, error) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &ProtoClient{Conn: conn}, nil
}

func (ps *ProtoServer) Broadcast(message any) error {

	stream := pb.Event{}

	return nil
}

func SendMessage(stream pb.Service_SendMessageServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			// Client has closed the stream.
			return nil
		}
		if err != nil {
			return err
		}

		// Process the incoming message

		// Simulate sending events back to the client
		events := []*pb.Event{
			{Event: "Event 1", Data: "Data 1"},
			{Event: "Event 2", Data: "Data 2"},
		}

		for _, event := range events {
			if err := stream.Send(&pb.Message{Message: }); err != nil {
				return err
			}
		}
	}
}

func (pc *ProtoClient) Receive() {
	for {

	}
}
