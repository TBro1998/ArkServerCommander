package docker_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// StatsSample 容器某一时刻的资源使用情况。CPU 是累计值，需要和上一次采样相减才能算出占用率
type StatsSample struct {
	CPUTotal    uint64
	SystemCPU   uint64
	OnlineCPUs  uint32
	MemoryUsage uint64
	MemoryLimit uint64
	Time        time.Time
}

// ContainerStats 获取容器的一次资源采样（不等待两次采样，立即返回）
func (dm *DockerManager) ContainerStats(containerName string) (*StatsSample, error) {
	resp, err := dm.client.ContainerStatsOneShot(dm.ctx, containerName)
	if err != nil {
		return nil, fmt.Errorf("获取容器资源占用失败: %v", err)
	}
	defer resp.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("解析容器资源占用失败: %v", err)
	}

	// 与 docker stats 一致：内存占用扣除可回收的页缓存（cgroup v2 为 inactive_file，v1 为 total_inactive_file）
	mem := stats.MemoryStats.Usage
	if v, ok := stats.MemoryStats.Stats["inactive_file"]; ok && v < mem {
		mem -= v
	} else if v, ok := stats.MemoryStats.Stats["total_inactive_file"]; ok && v < mem {
		mem -= v
	}

	cpus := stats.CPUStats.OnlineCPUs
	if cpus == 0 {
		cpus = uint32(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	return &StatsSample{
		CPUTotal:    stats.CPUStats.CPUUsage.TotalUsage,
		SystemCPU:   stats.CPUStats.SystemUsage,
		OnlineCPUs:  cpus,
		MemoryUsage: mem,
		MemoryLimit: stats.MemoryStats.Limit,
		Time:        stats.Read,
	}, nil
}

// CPUPercent 根据前后两次采样计算 CPU 占用率（100% 表示占满一个核心，与 docker stats 一致）
func CPUPercent(prev, cur *StatsSample) float64 {
	if prev == nil || cur == nil || cur.CPUTotal < prev.CPUTotal || cur.SystemCPU <= prev.SystemCPU {
		return 0
	}
	cpuDelta := float64(cur.CPUTotal - prev.CPUTotal)
	sysDelta := float64(cur.SystemCPU - prev.SystemCPU)
	cpus := float64(cur.OnlineCPUs)
	if cpus == 0 {
		cpus = 1
	}
	return cpuDelta / sysDelta * cpus * 100
}

// VolumeSizes 返回所有卷占用的磁盘空间（字节），key 是卷名。
// Docker 需要遍历卷内文件计算大小，比较慢，调用方应缓存结果
func (dm *DockerManager) VolumeSizes() (map[string]int64, error) {
	usage, err := dm.client.DiskUsage(dm.ctx, types.DiskUsageOptions{Types: []types.DiskUsageObject{types.VolumeObject}})
	if err != nil {
		return nil, fmt.Errorf("获取卷磁盘占用失败: %v", err)
	}
	result := make(map[string]int64, len(usage.Volumes))
	for _, v := range usage.Volumes {
		if v == nil {
			continue
		}
		size := int64(0)
		if v.UsageData != nil && v.UsageData.Size > 0 {
			size = v.UsageData.Size
		}
		result[v.Name] = size
	}
	return result, nil
}
