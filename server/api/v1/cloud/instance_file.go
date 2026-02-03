package cloud

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListContainerFiles 列出容器文件
func (instApi *InstanceApi) ListContainerFiles(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	containerID := c.Query("containerId")
	path := c.Query("path")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)

	if nodeID == 0 || containerID == "" {
		response.FailWithMessage("参数错误", c)
		return
	}

	files, err := instService.ListContainerFiles(c, nodeID, containerID, path)
	if err != nil {
		global.GVA_LOG.Error("获取文件列表失败", zap.Error(err))
		response.FailWithMessage("获取文件列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(files, "获取成功", c)
}

// DownloadContainerFile 下载容器文件
func (instApi *InstanceApi) DownloadContainerFile(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	containerID := c.Query("containerId")
	filePath := c.Query("path")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)

	if nodeID == 0 || containerID == "" || filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	reader, name, err := instService.DownloadContainerFile(c, nodeID, containerID, filePath)
	if err != nil {
		global.GVA_LOG.Error("下载文件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar\"", name))
	c.Header("Content-Type", "application/x-tar")
	
	io.Copy(c.Writer, reader)
}

// UploadContainerFile 上传文件到容器
func (instApi *InstanceApi) UploadContainerFile(c *gin.Context) {
	nodeIDStr := c.PostForm("nodeId")
	containerID := c.PostForm("containerId")
	destPath := c.PostForm("path")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("获取上传文件失败", c)
		return
	}
	defer file.Close()

	if nodeID == 0 || containerID == "" || destPath == "" {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = instService.UploadContainerFile(c, nodeID, containerID, destPath, file, header.Filename)
	if err != nil {
		global.GVA_LOG.Error("上传文件失败", zap.Error(err))
		response.FailWithMessage("上传文件失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("上传成功", c)
}

// RemoveContainerFile 删除容器文件
func (instApi *InstanceApi) RemoveContainerFile(c *gin.Context) {
	var req struct {
		NodeID      int64  `json:"nodeId"`
		ContainerID string `json:"containerId"`
		Path        string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err := instService.RemoveContainerFile(c, req.NodeID, req.ContainerID, req.Path)
	if err != nil {
		global.GVA_LOG.Error("删除文件失败", zap.Error(err))
		response.FailWithMessage("删除文件失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
