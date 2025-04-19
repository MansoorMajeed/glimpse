package server

import (
	"net"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
	"google.golang.org/grpc"
)

type GlimpseServer struct {
	pb.UnimplementedGlimpseServiceServer
}

// NewGlimpseServer creates a new instance of GlimpseServer
func NewGlimpseServer() *GlimpseServer {
	return &GlimpseServer{}
}

func StartGRPCServer() {
	logger.Info("Starting gRPC server on port 5001...")
	lis, err := net.Listen("tcp", ":5001")

	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGlimpseServiceServer(grpcServer, NewGlimpseServer())

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
