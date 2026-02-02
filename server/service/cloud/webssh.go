package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源，生产环境应该限制
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type string `json:"type"` // resize, input, ping
	Cols uint16 `json:"cols"` // 终端列数
	Rows uint16 `json:"rows"` // 终端行数
	Data string `json:"data"` // 输入数据
}

// ContainerSSHSession 容器 SSH 会话
type ContainerSSHSession struct {
	ContainerID  string
	ExecID       string
	Conn         *websocket.Conn
	HijackedResp types.HijackedResponse
	InputStream  io.WriteCloser
	OutputStream io.Reader
	Ctx          context.Context
	Cancel       context.CancelFunc
	DockerCli    *client.Client
	CreatedAt    time.Time
	LastActive   time.Time
}

// 会话管理器 - 用于跟踪活跃的 SSH 会话
var (
	sessionMutex   sync.RWMutex
	activeSessions = make(map[string]*ContainerSSHSession)
)

// RegisterSession 注册会话
func registerSession(sessionID string, session *ContainerSSHSession) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	session.CreatedAt = time.Now()
	session.LastActive = time.Now()
	activeSessions[sessionID] = session
	global.GVA_LOG.Info("SSH 会话已注册", zap.String("sessionID", sessionID))
}

// UnregisterSession 注销会话
func unregisterSession(sessionID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	if session, exists := activeSessions[sessionID]; exists {
		session.Cancel()
		delete(activeSessions, sessionID)
		global.GVA_LOG.Info("SSH 会话已注销", zap.String("sessionID", sessionID))
	}
}

// GetActiveSessions 获取活跃会话列表
func GetActiveSessions() []map[string]interface{} {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	var sessions []map[string]interface{}
	for id, sess := range activeSessions {
		sessions = append(sessions, map[string]interface{}{
			"sessionID":   id,
			"containerID": sess.ContainerID,
			"createdAt":   sess.CreatedAt,
			"lastActive":  sess.LastActive,
		})
	}
	return sessions
}

// UpdateSessionActivity 更新会话活动时间
func updateSessionActivity(sessionID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	if session, exists := activeSessions[sessionID]; exists {
		session.LastActive = time.Now()
	}
}

// cleanupInactiveSessions 清理不活跃的会话（超过30分钟无活动）
func cleanupInactiveSessions() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	now := time.Now()
	timeout := 30 * time.Minute

	for id, session := range activeSessions {
		if now.Sub(session.LastActive) > timeout {
			global.GVA_LOG.Info("清理不活跃会话", zap.String("sessionID", id))
			session.Cancel()
			session.Conn.Close()
			delete(activeSessions, id)
		}
	}
}

// WebSSHHandle 处理 WebSSH WebSocket 连接
func (instService *InstanceService) WebSSHHandle(instanceID string, conn *websocket.Conn) error {
	global.GVA_LOG.Info("WebSSH 连接请求", zap.String("instanceID", instanceID))

	// 生成会话ID
	sessionID := fmt.Sprintf("%s-%d", instanceID, time.Now().UnixNano())

	// 1. 获取实例信息
	var inst cloud.Instance
	if err := global.GVA_DB.First(&inst, instanceID).Error; err != nil {
		global.GVA_LOG.Error("获取实例信息失败", zap.Error(err))
		return fmt.Errorf("获取实例信息失败: %v", err)
	}

	// 2. 检查容器状态 (仅作日志记录，不强制阻断，交由 Docker API 判断)
	if inst.ContainerStatus != nil && *inst.ContainerStatus != "运行中" {
		global.GVA_LOG.Warn("尝试连接非运行状态容器",
			zap.String("status", *inst.ContainerStatus),
			zap.String("instanceID", instanceID))
	}

	if inst.DockerContainer == nil || *inst.DockerContainer == "" {
		return fmt.Errorf("容器ID不存在")
	}

	// 3. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, inst.NodeID).Error; err != nil {
		return fmt.Errorf("获取节点信息失败: %v", err)
	}

	// 4. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}

	containerID := *inst.DockerContainer

	// 5. 创建上下文和会话
	ctx, cancel := context.WithCancel(context.Background())

	// 6. 创建 Exec 实例 - 先尝试 bash
	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
		Env:          []string{"TERM=xterm-256color", "COLUMNS=80", "LINES=24"},
	}

	execResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		// 如果 bash 不存在，尝试 sh
		execConfig.Cmd = []string{"/bin/sh"}
		execResp, err = cli.ContainerExecCreate(ctx, containerID, execConfig)
		if err != nil {
			global.GVA_LOG.Error("创建 Exec 失败", zap.Error(err))
			return fmt.Errorf("创建 Exec 失败: %v", err)
		}
	}

	global.GVA_LOG.Info("Exec 创建成功",
		zap.String("execID", execResp.ID),
		zap.String("sessionID", sessionID))

	// 7. 连接到 Exec
	hijackedResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		global.GVA_LOG.Error("连接 Exec 失败", zap.Error(err))
		return fmt.Errorf("连接 Exec 失败: %v", err)
	}

	// 8. 启动 Exec
	if err := cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{
		Tty: true,
	}); err != nil {
		hijackedResp.Close()
		global.GVA_LOG.Error("启动 Exec 失败", zap.Error(err))
		return fmt.Errorf("启动 Exec 失败: %v", err)
	}

	// 9. 创建会话
	session := &ContainerSSHSession{
		ContainerID:  containerID,
		ExecID:       execResp.ID,
		Conn:         conn,
		HijackedResp: hijackedResp,
		InputStream:  hijackedResp.Conn,
		OutputStream: hijackedResp.Reader,
		Ctx:          ctx,
		Cancel:       cancel,
		DockerCli:    cli,
		CreatedAt:    time.Now(),
		LastActive:   time.Now(),
	}

	// 注册会话
	registerSession(sessionID, session)
	defer unregisterSession(sessionID)

	// 10. 发送欢迎消息
	welcomeMsg := fmt.Sprintf("\r\n\x1b[1;32mLoRAForge Web Terminal\x1b[0m\r\n")
	welcomeMsg += fmt.Sprintf("容器: %s\r\n", safeString(inst.InstanceName))
	welcomeMsg += fmt.Sprintf("节点: %s\r\n", safeString(node.Name))
	welcomeMsg += fmt.Sprintf("会话: %s\r\n", sessionID[:12])
	welcomeMsg += "\x1b[2m连接时间: \x1b[0m" + time.Now().Format("2006-01-02 15:04:05") + "\r\n\n"
	conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))

	// 11. 处理会话（包括消息处理和双向转发）
	return instService.handleSessionWithMessage(session)
}

// handleSessionWithMessage 处理 WebSocket 会话（带消息处理）
func (instService *InstanceService) handleSessionWithMessage(session *ContainerSSHSession) error {
	var wg sync.WaitGroup
	wg.Add(2) // 输入/消息处理、输出

	// 消息处理 goroutine
	go func() {
		defer wg.Done()
		instService.handleMessages(session)
	}()

	// 容器输出 -> WebSocket
	go func() {
		defer wg.Done()
		instService.forwardOutput(session)
	}()

	// WebSocket 输入 -> 容器
	// 注意：输入处理已合并到 handleMessages 中，避免并发读取 WebSocket 导致的竞争问题
	// go func() {
	// 	defer wg.Done()
	// 	instService.forwardInput(session)
	// }()

	// 等待任一方向结束
	wg.Wait()

	global.GVA_LOG.Info("WebSSH 会话结束",
		zap.String("containerID", session.ContainerID),
		zap.String("execID", session.ExecID))

	return nil
}

// handleMessages 处理 WebSocket 消息（resize、ping 等）
func (instService *InstanceService) handleMessages(session *ContainerSSHSession) {
	defer session.Conn.Close()

	// 设置 ping 间隔
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-session.Ctx.Done():
			return

		case <-pingTicker.C:
			// 发送心跳
			if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				global.GVA_LOG.Debug("发送心跳失败，连接可能已断开", zap.Error(err))
				return
			}

		default:
			// 读取消息
			_, message, err := session.Conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					global.GVA_LOG.Debug("读取 WebSocket 消息失败", zap.Error(err))
				}
				return
			}

			// 解析消息
			var wsMsg WSMessage
			// 尝试解析为 JSON
			if err := json.Unmarshal(message, &wsMsg); err == nil && (wsMsg.Type == "resize" || wsMsg.Type == "ping") {
				// JSON 消息（resize, ping 等）
				switch wsMsg.Type {
				case "resize":
					instService.handleResize(session, wsMsg.Rows, wsMsg.Cols)
				case "ping":
					// 响应 pong
					session.Conn.WriteMessage(websocket.PongMessage, nil)
					updateSessionActivity(fmt.Sprintf("%s-%s", session.ContainerID, session.ExecID))
				}
			} else {
				// 普通输入数据（非 JSON），直接转发到容器
				if _, err := session.InputStream.Write(message); err != nil {
					global.GVA_LOG.Error("写入容器输入失败", zap.Error(err))
					return
				}
				updateSessionActivity(fmt.Sprintf("%s-%s", session.ContainerID, session.ExecID))
			}
		}
	}
}

// handleResize 处理终端大小调整
func (instService *InstanceService) handleResize(session *ContainerSSHSession, rows, cols uint16) {
	if rows == 0 || cols == 0 {
		return
	}

	global.GVA_LOG.Debug("调整终端大小",
		zap.String("execID", session.ExecID),
		zap.Uint16("rows", rows),
		zap.Uint16("cols", cols))

	if err := session.DockerCli.ContainerExecResize(session.Ctx, session.ExecID, container.ResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	}); err != nil {
		global.GVA_LOG.Error("调整终端大小失败", zap.Error(err))
	}
}

// forwardOutput 将容器输出转发到 WebSocket
func (instService *InstanceService) forwardOutput(session *ContainerSSHSession) {
	buf := make([]byte, 8192)

	for {
		select {
		case <-session.Ctx.Done():
			return
		default:
			n, err := session.OutputStream.Read(buf)
			if err != nil {
				if err != io.EOF {
					global.GVA_LOG.Debug("读取容器输出失败", zap.Error(err))
				}
				return
			}

			if n > 0 {
				// 发送数据到 WebSocket
				if err := session.Conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					global.GVA_LOG.Debug("发送 WebSocket 消息失败", zap.Error(err))
					return
				}
			}

			// 更新活动时间
			updateSessionActivity(fmt.Sprintf("%s-%s", session.ContainerID, session.ExecID))
		}
	}
}

// forwardInput 将 WebSocket 输入转发到容器（优化版）
func (instService *InstanceService) forwardInput(session *ContainerSSHSession) {
	for {
		select {
		case <-session.Ctx.Done():
			return
		default:
			_, message, err := session.Conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					global.GVA_LOG.Debug("读取 WebSocket 输入失败", zap.Error(err))
				}
				return
			}

			// 写入到容器 stdin
			if _, err := session.InputStream.Write(message); err != nil {
				global.GVA_LOG.Error("写入容器输入失败", zap.Error(err))
				return
			}

			updateSessionActivity(fmt.Sprintf("%s-%s", session.ContainerID, session.ExecID))
		}
	}
}

// ResizeTty 调整终端大小
func (instService *InstanceService) ResizeTty(execID string, cli *client.Client, ctx context.Context, rows, cols uint16) error {
	return cli.ContainerExecResize(ctx, execID, container.ResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
}

// safeString 安全获取字符串值
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StartCleanupTimer 启动定期清理不活跃会话的定时器
func StartCleanupTimer() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			cleanupInactiveSessions()
		}
	}()
}
