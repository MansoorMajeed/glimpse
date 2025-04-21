package agentid

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mansoormajeed/glimpse/internal/common/logger"
)

func getAgentIDPath() string {
	if xdgState := os.Getenv("XDG_STATE_HOME"); xdgState != "" {
		return filepath.Join(xdgState, "glimpse", "agent-id")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalf("Unable to get home directory: %v", err)
	}

	return filepath.Join(homeDir, ".glimpse", "agent-id")
}

func LoadOrGenerateAgentID() string {
	path := getAgentIDPath()

	// Try to read existing ID
	data, err := os.ReadFile(path)
	if err == nil {
		return strings.TrimSpace(string(data))
	}

	// Generate new ID
	id := uuid.NewString()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		logger.Fatalf("Unable to create agent ID directory: %v", err)
	}

	// Save it
	if err := os.WriteFile(path, []byte(id), 0600); err != nil {
		logger.Fatalf("Unable to write agent ID file: %v", err)
	}

	return id
}
