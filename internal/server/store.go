package server

import (
	"sync"
	"time"

	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
)

type MetricEntry struct {
	Timestamp time.Time
	Metrics   *pb.AgentMetrics
}

type AgentData struct {
	Hostname     string
	OS           string
	LastSeen     time.Time
	ConnectedFor time.Duration

	MetricsHistory []MetricEntry
	metricsIndex   int // points to the next write position.
	metricsCount   int // tracks how many valid entries exist.
}

type ServerStore struct {
	sync.Mutex
	agents     map[string]*AgentData
	bufferSize int
}

func NewServerStore(bufferSize int) *ServerStore {
	return &ServerStore{
		agents:     make(map[string]*AgentData),
		bufferSize: bufferSize,
	}
}
