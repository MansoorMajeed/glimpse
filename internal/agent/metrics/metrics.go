package metrics

import (
	"time"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Metrics struct {
	CPUUsage        int64
	MemoryUsage     int64
	DiskUsage       int64
	NetworkUpload   int64
	NetworkDownload int64
	DiskReadKB      int64
	DiskWriteKB     int64
	CPUTemp         int64
	Uptime          int64 // Uptime in seconds
}

type AgentHeartbeat struct {
	Hostname     string
	OS           string
	Metrics      Metrics
	LastSeen     int64 // Last seen in seconds
	ConnectedFor int64 // Connected for in seconds
}

func GetAgentMetrics() (Metrics, error) {

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
		DiskReadKB:      0, // Placeholder for disk read KB
		DiskWriteKB:     0, // Placeholder for disk write KB
		CPUTemp:         cpuTemp,
	}
	return metrics, nil
}

func GetHostUptime() int64 {
	uptime, err := host.Uptime()
	if err != nil {
		logger.Errorf("Error getting uptime: %v", err)
		return 0
	}
	return int64(uptime)
}

func GetCPUUsage() int64 {
	cpuPercent, err := cpu.Percent(time.Second, false) // false means per cpu core? TODO: check
	if err != nil {
		logger.Errorf("Error getting cpu usage: %v", err)
		return 0
	}
	return int64(cpuPercent[0])
}

func GetMemoryUsage() int64 {
	memory, err := mem.VirtualMemory()
	if err != nil {
		logger.Errorf("Error getting memory usage: %v", err)
		return 0
	}
	return int64(memory.UsedPercent)
}

func GetDiskUsage() int64 {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		logger.Errorf("Error getting disk usage: %v", err)
		return 0
	}
	return int64(diskUsage.UsedPercent)
}
func GetNetworkUsage() (int64, int64) {
	// Placeholder values for network upload and download
	// In a real implementation, you would use a library to get actual network stats
	upload := int64(0)
	download := int64(0)

	return upload, download
}

func GetCPUTemperature() int64 {
	// Placeholder value for CPU temperature
	// In a real implementation, you would use a library to get actual CPU temperature
	temp := int64(0)

	return temp
}
