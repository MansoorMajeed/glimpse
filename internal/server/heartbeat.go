package server

import (
	"context"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
)

func (s *GlimpseServer) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	// Log the heartbeat request
	logger.Debugf("Received heartbeat : %v", req)

	resp := &pb.HeartbeatResponse{
		Message:      "Heartbeat received",
		Success:      true,
		StatusCode:   200,
		ErrorMessage: "",
	}

	return resp, nil
}
