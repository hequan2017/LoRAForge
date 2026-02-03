package cloud

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VolumeApi struct{}

var volumeService = service.ServiceGroupApp.CloudServiceGroup.VolumeService

// GetVolumes 获取卷列表
func (api *VolumeApi) GetVolumes(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)
	if nodeID == 0 {
		response.FailWithMessage("节点ID不能为空", c)
		return
	}

	list, err := volumeService.GetVolumes(c, nodeID)
	if err != nil {
		global.GVA_LOG.Error("获取卷列表失败", zap.Error(err))
		response.FailWithMessage("获取卷列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(list.Volumes, "获取成功", c)
}

// CreateVolume 创建卷
func (api *VolumeApi) CreateVolume(c *gin.Context) {
	var req struct {
		NodeID int64  `json:"nodeId"`
		Name   string `json:"name"`
		Driver string `json:"driver"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if req.NodeID == 0 || req.Name == "" {
		response.FailWithMessage("节点ID和卷名称不能为空", c)
		return
	}

	vol, err := volumeService.CreateVolume(c, req.NodeID, req.Name, req.Driver)
	if err != nil {
		global.GVA_LOG.Error("创建卷失败", zap.Error(err))
		response.FailWithMessage("创建卷失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(vol, "创建成功", c)
}

// RemoveVolume 删除卷
func (api *VolumeApi) RemoveVolume(c *gin.Context) {
	var req struct {
		NodeID   int64  `json:"nodeId"`
		VolumeID string `json:"volumeId"`
		Force    bool   `json:"force"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if req.NodeID == 0 || req.VolumeID == "" {
		response.FailWithMessage("节点ID和卷ID不能为空", c)
		return
	}

	err := volumeService.RemoveVolume(c, req.NodeID, req.VolumeID, req.Force)
	if err != nil {
		global.GVA_LOG.Error("删除卷失败", zap.Error(err))
		response.FailWithMessage("删除卷失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(nil, "删除成功", c)
}
