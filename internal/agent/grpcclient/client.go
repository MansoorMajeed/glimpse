package grpcclient

import (
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
	"google.golang.org/grpc"
)

type GlimpseServiceClient struct {
	client pb.GlimpseServiceClient
}

func NewGlimpseServiceClient() (*GlimpseServiceClient, error) {

	var opts []grpc.DialOption

	conn, err := grpc.NewClient("localhost:5001", opts...)
	if err != nil {
		return nil, err
	}

	c := pb.NewGlimpseServiceClient(conn)
	return &GlimpseServiceClient{
		client: c,
	}, nil
}
