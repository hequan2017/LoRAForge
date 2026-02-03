package cloud

import (
	"bufio"
	"io"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct{}

var imageService = service.ServiceGroupApp.CloudServiceGroup.ImageService

// GetImages 获取镜像列表
func (api *ImageApi) GetImages(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)
	if nodeID == 0 {
		response.FailWithMessage("节点ID不能为空", c)
		return
	}

	list, err := imageService.GetImages(c, nodeID)
	if err != nil {
		global.GVA_LOG.Error("获取镜像列表失败", zap.Error(err))
		response.FailWithMessage("获取镜像列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// PullImage 拉取镜像
// 此接口返回流式数据
func (api *ImageApi) PullImage(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	imageName := c.Query("imageName")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)

	if nodeID == 0 || imageName == "" {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	reader, err := imageService.PullImage(c, nodeID, imageName)
	if err != nil {
		global.GVA_LOG.Error("拉取镜像失败", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 逐行读取并返回
	scanner := bufio.NewScanner(reader)
	c.Stream(func(w io.Writer) bool {
		if scanner.Scan() {
			c.Writer.Write([]byte("data: "))
			c.Writer.Write(scanner.Bytes())
			c.Writer.Write([]byte("\n\n"))
			c.Writer.Flush()
			return true
		}
		return false
	})
}

// RemoveImage 删除镜像
func (api *ImageApi) RemoveImage(c *gin.Context) {
	var req struct {
		NodeID  int64  `json:"nodeId"`
		ImageID string `json:"imageId"`
		Force   bool   `json:"force"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	res, err := imageService.RemoveImage(c, req.NodeID, req.ImageID, req.Force)
	if err != nil {
		global.GVA_LOG.Error("删除镜像失败", zap.Error(err))
		response.FailWithMessage("删除镜像失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(res, "删除成功", c)
}
