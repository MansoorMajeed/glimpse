package server

import (
	"sync"
	"time"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	pb "github.com/mansoormajeed/glimpse/pkg/pb/proto"
)

type MetricEntry struct {
	Timestamp time.Time
	Metrics   *pb.AgentMetrics
}

type AgentData struct {
	AgentID      string
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

func (s *ServerStore) AddOrUpdateAgent(req *pb.HeartbeatRequest) {
	s.Lock()
	defer s.Unlock()
	logger.Debugf("Adding/updating agent: %s", req.Hostname)
	agent, exists := s.agents[req.Hostname]
	if !exists {
		logger.Debugf("Creating new agent entry for: %s", req.Hostname)
		agent = &AgentData{
			AgentID:        req.AgentId,
			Hostname:       req.Hostname,
			OS:             req.Os,
			MetricsHistory: make([]MetricEntry, s.bufferSize),
		}
		s.agents[req.Hostname] = agent
	}

	agent.LastSeen = time.Now()
	agent.ConnectedFor = time.Duration(req.ConnectedFor) * time.Second

	entry := MetricEntry{
		Timestamp: time.Now(),
		Metrics:   req.Metrics,
	}
	agent.MetricsHistory[agent.metricsIndex] = entry

	agent.metricsIndex = (agent.metricsIndex + 1) % s.bufferSize
	if agent.metricsCount < s.bufferSize {
		agent.metricsCount++
	}

	logger.Debugf("added to the index: %d", agent.metricsIndex)
}

func (s *ServerStore) GetAgentData(hostname string) (*AgentData, bool) {
	s.Lock()
	defer s.Unlock()

	agent, exists := s.agents[hostname]
	if !exists {
		return nil, false
	}

	// Return a copy of the agent data to prevent external modification
	agentCopy := *agent
	return &agentCopy, true
}
