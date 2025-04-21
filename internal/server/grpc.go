package server

import (
	"context"
	"net"
	"time"

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
func NewGlimpseServer() *GlimpseServer {
	return &GlimpseServer{
		store: NewServerStore(60), // Set the buffer size to 60 for 1 minute
	}
}

func (s *GlimpseServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	// Log the heartbeat request
	logger.Debugf("Received heartbeat : %v", req)

	s.store.Lock()
	defer s.store.Unlock()

	agent, exists := s.store.agents[req.Hostname]
	if !exists {
		agent = &AgentData{
			Hostname:       req.Hostname,
			OS:             req.Os,
			MetricsHistory: make([]MetricEntry, s.store.bufferSize),
		}
		s.store.agents[req.Hostname] = agent
	}

	agent.LastSeen = time.Now()
	agent.ConnectedFor = time.Duration(req.ConnectedFor) * time.Second

	// add to the ring buffer
	entry := MetricEntry{
		Timestamp: time.Now(),
		Metrics:   req.Metrics,
	}
	agent.MetricsHistory[agent.metricsIndex] = entry
	agent.metricsIndex = (agent.metricsIndex + 1) % s.store.bufferSize
	if agent.metricsCount < s.store.bufferSize {
		agent.metricsCount++
	}

	resp := &pb.HeartbeatResponse{
		Message:      "Heartbeat received",
		Success:      true,
		StatusCode:   200,
		ErrorMessage: "",
	}
	logger.Debugf(util.PrettyYaml(s.store.agents))
	return resp, nil
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
