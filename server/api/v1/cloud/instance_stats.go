package cloud

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/cloud"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetContainerStats 获取容器实时监控数据 (SSE)
func (instApi *InstanceApi) GetContainerStats(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	containerID := c.Query("containerId")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)

	if nodeID == 0 || containerID == "" {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	reader, err := instService.GetContainerStats(c, nodeID, containerID)
	if err != nil {
		global.GVA_LOG.Error("获取容器监控失败", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	scanner := bufio.NewScanner(reader)
	c.Stream(func(w io.Writer) bool {
		if !scanner.Scan() {
			return false
		}
		line := scanner.Bytes()
		// 解析 JSON 并计算 CPU 使用率等简化数据发给前端
		var stats cloud.ContainerStatsJSON
		if err := json.Unmarshal(line, &stats); err == nil {
			cpuPercent := cloud.CalculateCPUUsage(stats)
			memUsage := stats.MemoryStats.Usage
			memLimit := stats.MemoryStats.Limit
			memPercent := 0.0
			if memLimit > 0 {
				memPercent = float64(memUsage) / float64(memLimit) * 100.0
			}

			netRx := uint64(0)
			netTx := uint64(0)
			for _, net := range stats.Networks {
				netRx += net.RxBytes
				netTx += net.TxBytes
			}

			// 计算磁盘IO
			blkRead := uint64(0)
			blkWrite := uint64(0)
			for _, io := range stats.BlkioStats.IoServiceBytesRecursive {
				if io.Op == "Read" {
					blkRead += io.Value
				} else if io.Op == "Write" {
					blkWrite += io.Value
				}
			}

			// 构建简化版数据
			simpleStats := gin.H{
				"cpu_percent": cpuPercent,
				"mem_usage":   memUsage,
				"mem_limit":   memLimit,
				"mem_percent": memPercent,
				"net_rx":      netRx,
				"net_tx":      netTx,
				"blk_read":    blkRead,
				"blk_write":   blkWrite,
				"timestamp":   time.Now().Unix(),
			}

			data, _ := json.Marshal(simpleStats)
			c.Writer.Write([]byte("data: "))
			c.Writer.Write(data)
			c.Writer.Write([]byte("\n\n"))
			c.Writer.Flush()
		}
		return true
	})
}
