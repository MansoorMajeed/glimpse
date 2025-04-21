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
	logger.Debugf("Adding/updating agent: %s", req.AgentId)
	agent, exists := s.agents[req.AgentId]
	if !exists {
		logger.Debugf("Creating new agent entry for: %s", req.Hostname)
		agent = &AgentData{
			AgentID:        req.AgentId,
			Hostname:       req.Hostname,
			OS:             req.Os,
			MetricsHistory: make([]MetricEntry, s.bufferSize),
		}
		s.agents[req.AgentId] = agent
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

func (s *ServerStore) GetAllAgents() []*AgentData {
	s.Lock()
	defer s.Unlock()

	agentList := make([]*AgentData, 0, len(s.agents))
	for _, agent := range s.agents {
		// Return a copy of the agent data to prevent external modification
		agentCopy := *agent
		agentList = append(agentList, &agentCopy)
	}
	return agentList
}

func (s *ServerStore) GetAgentData(agentId string) (*AgentData, bool) {
	s.Lock()
	defer s.Unlock()

	agent, exists := s.agents[agentId]
	if !exists {
		return nil, false
	}

	// Return a copy of the agent data to prevent external modification
	agentCopy := *agent
	return &agentCopy, true
}

func (a *AgentData) Latest() *pb.AgentMetrics {
	if a.metricsCount == 0 {
		return nil
	}
	idx := (a.metricsIndex - 1 + len(a.MetricsHistory)) % len(a.MetricsHistory)
	return a.MetricsHistory[idx].Metrics
}
