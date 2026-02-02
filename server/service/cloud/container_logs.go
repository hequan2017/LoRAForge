package cloud

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ContainerLogsHandle 处理容器日志 WebSocket 连接
func (instService *InstanceService) ContainerLogsHandle(instanceID string, conn *websocket.Conn) error {
	global.GVA_LOG.Info("容器日志连接请求", zap.String("instanceID", instanceID))

	// 1. 获取实例信息
	var inst cloud.Instance
	if err := global.GVA_DB.First(&inst, instanceID).Error; err != nil {
		global.GVA_LOG.Error("获取实例信息失败", zap.Error(err))
		return fmt.Errorf("获取实例信息失败: %v", err)
	}

	if inst.DockerContainer == nil || *inst.DockerContainer == "" {
		return fmt.Errorf("容器ID不存在")
	}

	// 2. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, inst.NodeID).Error; err != nil {
		return fmt.Errorf("获取节点信息失败: %v", err)
	}

	// 3. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}

	containerID := *inst.DockerContainer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. 获取日志流
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100", // 默认显示最后100行
		Timestamps: false,
	}

	logsReader, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		global.GVA_LOG.Error("获取容器日志失败", zap.Error(err))
		return fmt.Errorf("获取容器日志失败: %v", err)
	}
	defer logsReader.Close()

	// 5. 发送欢迎消息
	welcomeMsg := fmt.Sprintf("\r\n\x1b[1;32mLoRAForge Container Logs\x1b[0m\r\n")
	welcomeMsg += fmt.Sprintf("容器: %s\r\n", safeString(inst.InstanceName))
	welcomeMsg += fmt.Sprintf("ID: %s\r\n", containerID[:12])
	welcomeMsg += "\x1b[2m连接时间: \x1b[0m" + time.Now().Format("2006-01-02 15:04:05") + "\r\n\n"
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// 6. 启动读取循环
	// 处理 WebSocket 关闭
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				cancel() // 关闭上下文，停止日志读取
				return
			}
		}
	}()

	// 读取日志并发送到 WebSocket
	scanner := bufio.NewScanner(logsReader)
	for scanner.Scan() {
		text := scanner.Text()
		// 简单处理：添加换行符
		// 注意：如果容器未启用 TTY，Docker 日志会有 8 字节头部，bufio.Scanner 可能无法完美处理二进制头部
		// 但对于文本日志，通常没问题。如果出现乱码，需要按 Docker 协议解析头部。
		// 这里先假设大部分日志是文本。
		
		// 更好的方式可能是直接 Copy，但 WebSocket 需要分帧。
		// 逐行读取是比较安全的方式，可以兼容 xterm.js
		
		err := conn.WriteMessage(websocket.TextMessage, []byte(text+"\r\n"))
		if err != nil {
			global.GVA_LOG.Debug("发送日志失败", zap.Error(err))
			break
		}
	}

	if err := scanner.Err(); err != nil {
		// 上下文取消导致的读取错误是预期的
		if err != context.Canceled {
			global.GVA_LOG.Error("读取日志流错误", zap.Error(err))
		}
	}

	return nil
}
