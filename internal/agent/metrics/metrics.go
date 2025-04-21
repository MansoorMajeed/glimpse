package metrics

import (
	"strings"
	"time"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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
	LastSeen     time.Time
	ConnectedFor time.Duration
}

var nwPreviousUpload, nwPreviousDownload uint64
var nwPreviousTime time.Time
var diskPrevReadBytes, diskPrevWriteBytes uint64
var diskPrevTime time.Time

func GetAgentMetrics() (Metrics, error) {

	cpuUsage := GetCPUUsage()
	memoryUsage := GetMemoryUsage()
	diskUsage := GetDiskUsage()
	networkUpload, networkDownload := GetNetworkUsage()
	cpuTemp := GetCPUTemperature()
	diskReadKB, diskWriteKB := GetDiskIO()
	uptime := GetHostUptime()

	metrics := Metrics{
		CPUUsage:        cpuUsage,
		MemoryUsage:     memoryUsage,
		DiskUsage:       diskUsage,
		NetworkUpload:   networkUpload,
		NetworkDownload: networkDownload,
		DiskReadKB:      diskReadKB,
		DiskWriteKB:     diskWriteKB,
		CPUTemp:         cpuTemp,
		Uptime:          uptime,
	}
	return metrics, nil
}

func GetHostUptime() int64 {
	uptime, err := host.Uptime()
	logger.Debugf("Uptime: %v", uptime)
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

	// sum the bytes sent and received across all network interfaces
	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		logger.Errorf("Error getting network usage: %v", err)
		return 0, 0
	}

	now := time.Now()
	duration := now.Sub(nwPreviousTime).Seconds()
	nwPreviousTime = now

	upload := counters[0].BytesSent
	download := counters[0].BytesRecv

	if nwPreviousUpload == 0 {
		nwPreviousUpload = upload
		nwPreviousDownload = download
		return 0, 0 // first call, no previous data
	}

	uploadRate := float64(upload-nwPreviousUpload) / duration
	downloadRate := float64(download-nwPreviousDownload) / duration

	nwPreviousUpload = upload
	nwPreviousDownload = download

	return int64(uploadRate / 1024), int64(downloadRate / 1024) // convert to KB
}

func GetDiskIO() (int64, int64) {
	ioStats, err := disk.IOCounters()
	if err != nil {
		logger.Errorf("Error getting disk IO: %v", err)
		return 0, 0
	}

	var readBytes, writeBytes uint64
	for _, stat := range ioStats {
		readBytes += stat.ReadBytes
		writeBytes += stat.WriteBytes
	}

	now := time.Now()
	delta := now.Sub(diskPrevTime).Seconds()
	diskPrevTime = now

	if diskPrevReadBytes == 0 {
		diskPrevReadBytes = readBytes
		diskPrevWriteBytes = writeBytes
		return 0, 0 // Skip first sample
	}

	readRate := float64(readBytes-diskPrevReadBytes) / delta
	writeRate := float64(writeBytes-diskPrevWriteBytes) / delta

	diskPrevReadBytes = readBytes
	diskPrevWriteBytes = writeBytes

	return int64(readRate / 1024), int64(writeRate / 1024) // in KB/s
}

// func GetCPUTemperature() int64 {

// 	temps, err := host.SensorsTemperatures()
// 	if err != nil || len(temps) == 0 {
// 		logger.Errorf("Error getting CPU temperature: %v", err)
// 		return 0
// 	}

// 	for _, t := range temps {
// 		if t.SensorKey == "Package id 0" || t.SensorKey == "Tctl" || t.SensorKey == "Core 0" {
// 			return int64(t.Temperature)
// 		}
// 	}

// 	// Fallback to the first temperature reading if specific sensor not found
// 	// probably wrong -- tomorrow's problem
// 	return int64(temps[0].Temperature)
// }

// This is tough. The sensor names are very inconsistent across hardware and platforms.
// So here's the approach I'm taking:
//  1. Try to find temperature sensors whose keys contain common CPU-related terms
//     like "cpu", "core", "package", or "tctl" (case-insensitive).
//  2. If multiple CPU-like sensors are found, pick the one with the highest temperature.
//     This gives a conservative estimate and avoids underreporting.
//  3. If no CPU-like sensors are found, fall back to the highest temperature across all sensors.
//     This isn't perfect, but avoids returning something irrelevant like NVMe or USB temps.
//  4. Long-term: this could be made configurable, or use smarter detection per hardware/vendor.
func GetCPUTemperature() int64 {
	temps, err := host.SensorsTemperatures()
	if err != nil || len(temps) == 0 {
		logger.Errorf("Error getting CPU temperature: %v", err)
		return 0
	}

	var bestTemp float64
	var found bool

	for _, t := range temps {
		key := strings.ToLower(t.SensorKey)
		if strings.Contains(key, "cpu") ||
			strings.Contains(key, "core") ||
			strings.Contains(key, "package") ||
			strings.Contains(key, "tctl") {
			if t.Temperature > bestTemp || !found {
				bestTemp = t.Temperature
				found = true
			}
		}
	}

	if found {
		return int64(bestTemp)
	}

	// fallback: max of all temps
	maxTemp := temps[0].Temperature
	for _, t := range temps {
		if t.Temperature > maxTemp {
			maxTemp = t.Temperature
		}
	}

	return int64(maxTemp)
}
