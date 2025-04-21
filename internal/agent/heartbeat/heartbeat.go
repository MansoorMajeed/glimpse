package heartbeat

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/mansoormajeed/glimpse/internal/agent/agentid"
	m "github.com/mansoormajeed/glimpse/internal/agent/metrics"
	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/mansoormajeed/glimpse/internal/common/logger/util"

	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
)

type HeartbeatService struct {
	client pb.GlimpseServiceClient
}

var agentID string

func NewHeartbeatService(client pb.GlimpseServiceClient) *HeartbeatService {
	return &HeartbeatService{
		client: client,
	}
}

func (h *HeartbeatService) Start(ctx context.Context) {

	logger.Info("Starting Heartbeat Service...")
	agentID = agentid.LoadOrGenerateAgentID()
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopping Heartbeat Service...")
				return
			case <-ticker.C:
				err := h.SendHeartbeat()
				if err != nil {
					logger.Errorf("Error sending heartbeat: %v", err)
				} else {
					logger.Info("Heartbeat sent successfully")
				}
			}
		}
	}()
	<-ctx.Done()
}

func (h *HeartbeatService) SendHeartbeat() error {

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	os := runtime.GOOS
	metrics, err := m.GetAgentMetrics()
	if err != nil {
		return fmt.Errorf("error getting agent metrics: %v", err)
	}
	req := &pb.HeartbeatRequest{
		Hostname: hostname,
		Os:       os,
		AgentId:  agentID,
		Metrics: &pb.AgentMetrics{
			CpuUsage:        metrics.CPUUsage,
			MemoryUsage:     metrics.MemoryUsage,
			DiskUsage:       metrics.DiskUsage,
			NetworkUpload:   metrics.NetworkUpload,
			NetworkDownload: metrics.NetworkDownload,
			DiskRead:        metrics.DiskReadKB,
			DiskWrite:       metrics.DiskWriteKB,
			CpuTemp:         metrics.CPUTemp,
			Uptime:          metrics.Uptime,
		},
	}
	logger.Info("Sending heartbeat request.... Hostname: ", hostname)
	logger.Debugf(util.PrettyYaml(req))

	resp, err := h.client.Heartbeat(context.Background(), req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("heartbeat failed: %s", resp.ErrorMessage)
	}

	return nil
}
