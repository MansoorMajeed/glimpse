package metrics

import (
	"os"
	"time"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Metrics struct {
	CPUUsage        int32
	MemoryUsage     int32
	DiskUsage       int32
	NetworkUpload   int32
	NetworkDownload int32
	CPUTemp         int32
}

type AgentHeartbeat struct {
	Hostname     string
	Metrics      Metrics
	LastSeen     time.Time
	ConnectedFor time.Duration
}

func GetAgentMetrics() (AgentHeartbeat, error) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Errorf("Error getting hostname: %v", err)
		return AgentHeartbeat{}, err
	}

	uptime := GetHostUptime()
	cpuUsage := GetCPUUsage()
	memoryUsage := GetMemoryUsage()
	diskUsage := GetDiskUsage()
	networkUpload, networkDownload := GetNetworkUsage()
	cpuTemp := GetCPUTemperature()

	metrics := Metrics{
		CPUUsage:        cpuUsage,
		MemoryUsage:     memoryUsage,
		DiskUsage:       diskUsage,
		NetworkUpload:   networkUpload,
		NetworkDownload: networkDownload,
		CPUTemp:         cpuTemp,
	}

	return AgentHeartbeat{
		Hostname:     hostname,
		Metrics:      metrics,
		LastSeen:     time.Now(),
		ConnectedFor: uptime,
	}, nil
}

func GetHostUptime() time.Duration {
	uptime, err := host.Uptime()
	if err != nil {
		logger.Errorf("Error getting uptime: %v", err)
		return 0
	}
	return time.Duration(uptime) * time.Second
}

func GetCPUUsage() int32 {
	cpuPercent, err := cpu.Percent(time.Second, false) // false means per cpu core? TODO: check
	if err != nil {
		logger.Errorf("Error getting cpu usage: %v", err)
		return 0
	}
	return int32(cpuPercent[0])
}

func GetMemoryUsage() int32 {
	memory, err := mem.VirtualMemory()
	if err != nil {
		logger.Errorf("Error getting memory usage: %v", err)
		return 0
	}
	return int32(memory.UsedPercent)
}

func GetDiskUsage() int32 {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		logger.Errorf("Error getting disk usage: %v", err)
		return 0
	}
	return int32(diskUsage.UsedPercent)
}
func GetNetworkUsage() (int32, int32) {
	// Placeholder values for network upload and download
	// In a real implementation, you would use a library to get actual network stats
	upload := int32(0)
	download := int32(0)

	return upload, download
}

func GetCPUTemperature() int32 {
	// Placeholder value for CPU temperature
	// In a real implementation, you would use a library to get actual CPU temperature
	temp := int32(0)

	return temp
}
