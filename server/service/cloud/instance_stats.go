package cloud

import (
	"context"
	"fmt"
	"io"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"go.uber.org/zap"
)

// GetContainerStats 获取容器实时监控数据 (流式)
func (instService *InstanceService) GetContainerStats(ctx context.Context, nodeID int64, containerID string) (io.ReadCloser, error) {
	// 1. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	// 2. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	// 注意：这里不能立即 Close，因为返回的是流，需要在调用方 Close 或者流结束时 Close
	// 但 cli.ContainerStats 返回的 body 关闭时并不会关闭 cli，所以我们需要一种机制来关闭 cli
	// 简单起见，我们在 stats 结束时依靠 GC 或显式关闭。
	// 更好的做法是封装一个 ReadCloser，同时关闭 cli。
	
	// 这里为了避免连接泄露，我们封装一下
	stats, err := cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		cli.Close()
		global.GVA_LOG.Error("获取容器监控失败", zap.Error(err))
		return nil, err
	}

	return &DockerStatsReader{
		Body: stats.Body,
		Cli:  cli,
	}, nil
}

// DockerStatsReader 封装 Docker Stats 流和 Client 关闭
type DockerStatsReader struct {
	Body io.ReadCloser
	Cli  io.Closer
}

func (r *DockerStatsReader) Read(p []byte) (n int, err error) {
	return r.Body.Read(p)
}

func (r *DockerStatsReader) Close() error {
	r.Body.Close()
	return r.Cli.Close()
}

// 辅助结构体：用于解析 Docker Stats JSON
type ContainerStatsJSON struct {
	Read      string `json:"read"`
	Preread   string `json:"preread"`
	PidsStats struct {
		Current int `json:"current"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []struct {
			Major uint64 `json:"major"`
			Minor uint64 `json:"minor"`
			Op    string `json:"op"`
			Value uint64 `json:"value"`
		} `json:"io_service_bytes_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        uint64   `json:"total_usage"`
			PercpuUsage       []uint64 `json:"percpu_usage"`
			UsageInKernelmode uint64   `json:"usage_in_kernelmode"`
			UsageInUsermode   uint64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCpus     int    `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        uint64   `json:"total_usage"`
			PercpuUsage       []uint64 `json:"percpu_usage"`
			UsageInKernelmode uint64   `json:"usage_in_kernelmode"`
			UsageInUsermode   uint64   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCpus     int    `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    uint64 `json:"usage"`
		MaxUsage uint64 `json:"max_usage"`
		Stats    struct {
			ActiveAnon              uint64 `json:"active_anon"`
			ActiveFile              uint64 `json:"active_file"`
			Cache                   uint64 `json:"cache"`
			Dirty                   uint64 `json:"dirty"`
			HierarchicalMemoryLimit uint64 `json:"hierarchical_memory_limit"`
			HierarchicalMemswLimit  uint64 `json:"hierarchical_memsw_limit"`
			InactiveAnon            uint64 `json:"inactive_anon"`
			InactiveFile            uint64 `json:"inactive_file"`
			MappedFile              uint64 `json:"mapped_file"`
			Pgfault                 uint64 `json:"pgfault"`
			Pgmajfault              uint64 `json:"pgmajfault"`
			Pgpgin                  uint64 `json:"pgpgin"`
			Pgpgout                 uint64 `json:"pgpgout"`
			Rss                     uint64 `json:"rss"`
			RssHuge                 uint64 `json:"rss_huge"`
			TotalActiveAnon         uint64 `json:"total_active_anon"`
			TotalActiveFile         uint64 `json:"total_active_file"`
			TotalCache              uint64 `json:"total_cache"`
			TotalDirty              uint64 `json:"total_dirty"`
			TotalInactiveAnon       uint64 `json:"total_inactive_anon"`
			TotalInactiveFile       uint64 `json:"total_inactive_file"`
			TotalMappedFile         uint64 `json:"total_mapped_file"`
			TotalPgfault            uint64 `json:"total_pgfault"`
			TotalPgmajfault         uint64 `json:"total_pgmajfault"`
			TotalPgpgin             uint64 `json:"total_pgpgin"`
			TotalPgpgout            uint64 `json:"total_pgpgout"`
			TotalRss                uint64 `json:"total_rss"`
			TotalRssHuge            uint64 `json:"total_rss_huge"`
			TotalUnevictable        uint64 `json:"total_unevictable"`
			TotalWriteback          uint64 `json:"total_writeback"`
			Unevictable             uint64 `json:"unevictable"`
			Writeback               uint64 `json:"writeback"`
		} `json:"stats"`
		Limit uint64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks map[string]struct {
		RxBytes   uint64 `json:"rx_bytes"`
		RxPackets uint64 `json:"rx_packets"`
		RxErrors  uint64 `json:"rx_errors"`
		RxDropped uint64 `json:"rx_dropped"`
		TxBytes   uint64 `json:"tx_bytes"`
		TxPackets uint64 `json:"tx_packets"`
		TxErrors  uint64 `json:"tx_errors"`
		TxDropped uint64 `json:"tx_dropped"`
	} `json:"networks"`
}

// CalculateCPUUsage 计算 CPU 使用率百分比
func CalculateCPUUsage(stats ContainerStatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PrecpuStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemCPUUsage) - float64(stats.PrecpuStats.SystemCPUUsage)
	onlineCPUs := float64(stats.CPUStats.OnlineCpus)
	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return 0.0
}
