package cloud

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ContainerLogs 容器日志 WebSocket 连接
// @Tags Instance
// @Summary 容器日志连接
// @Description 通过 WebSocket 获取容器日志
// @Security ApiKeyAuth
// @Param id query string true "实例ID"
// @Router /inst/logs [get]
func (instApi *InstanceApi) ContainerLogs(c *gin.Context) {
	// 获取实例ID
	instanceID := c.Query("id")
	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "实例ID不能为空"})
		return
	}

	// 升级为 WebSocket 连接
	// 注意：upgrader 在 webssh.go 中定义，同一个包下可以直接使用
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("WebSocket 升级失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "WebSocket 升级失败"})
		return
	}
	defer conn.Close()

	global.GVA_LOG.Info("ContainerLogs WebSocket 连接建立", zap.String("instanceID", instanceID))

	// 调用服务层处理
	err = service.ServiceGroupApp.CloudServiceGroup.InstanceService.ContainerLogsHandle(instanceID, conn)
	if err != nil {
		global.GVA_LOG.Error("ContainerLogs 处理失败", zap.Error(err))
		// 发送错误消息到客户端
		conn.WriteMessage(websocket.TextMessage, []byte("\r\n\x1b[31m错误: "+err.Error()+"\x1b[0m\r\n"))
	}
}
