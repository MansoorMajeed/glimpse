package grpcclient

import (
	"os"

	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
	"google.golang.org/grpc"
)

type GlimpseServiceClient struct {
	client pb.GlimpseServiceClient
}

func NewGlimpseServiceClient() (*GlimpseServiceClient, error) {
	var opts []grpc.DialOption

	// Get server address from environment variable or use default
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "localhost:5001" // Default fallback
	}

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	c := pb.NewGlimpseServiceClient(conn)
	return &GlimpseServiceClient{
		client: c,
	}, nil
}
