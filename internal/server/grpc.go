package server

import (
	"context"
	"net"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/mansoormajeed/glimpse/internal/common/logger/util"
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
	"google.golang.org/grpc"
)

type GlimpseServer struct {
	pb.UnimplementedGlimpseServiceServer
	store *ServerStore
}

// NewGlimpseServer creates a new instance of GlimpseServer
func NewGlimpseServer(store *ServerStore) *GlimpseServer {
	return &GlimpseServer{
		store: store, // for now use buffer 5. later change to 60
	}
}

func (s *GlimpseServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	// Log the heartbeat request
	logger.Debugf("Received heartbeat : %v", req)

	s.store.AddOrUpdateAgent(req)
	logger.Debugf("updated agent: %v", req.Hostname)

	resp := &pb.HeartbeatResponse{
		Message:      "Heartbeat received",
		Success:      true,
		StatusCode:   200,
		ErrorMessage: "",
	}
	logger.Debugf(util.PrettyYaml(s.store.agents))
	return resp, nil
}

func StartGRPCServer(store *ServerStore) {
	logger.Info("Starting gRPC server on port 5001...")
	lis, err := net.Listen("tcp", ":5001")

	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	glimpseServer := NewGlimpseServer(store)
	pb.RegisterGlimpseServiceServer(grpcServer, glimpseServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
