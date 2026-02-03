package cloud

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NetworkApi struct{}

var networkService = service.ServiceGroupApp.CloudServiceGroup.NetworkService

// GetNetworks 获取网络列表
func (api *NetworkApi) GetNetworks(c *gin.Context) {
	nodeIDStr := c.Query("nodeId")
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 64)
	if nodeID == 0 {
		response.FailWithMessage("节点ID不能为空", c)
		return
	}

	list, err := networkService.GetNetworks(c, nodeID)
	if err != nil {
		global.GVA_LOG.Error("获取网络列表失败", zap.Error(err))
		response.FailWithMessage("获取网络列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// CreateNetwork 创建网络
func (api *NetworkApi) CreateNetwork(c *gin.Context) {
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
		response.FailWithMessage("节点ID和网络名称不能为空", c)
		return
	}

	id, err := networkService.CreateNetwork(c, req.NodeID, req.Name, req.Driver)
	if err != nil {
		global.GVA_LOG.Error("创建网络失败", zap.Error(err))
		response.FailWithMessage("创建网络失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"id": id}, "创建成功", c)
}

// RemoveNetwork 删除网络
func (api *NetworkApi) RemoveNetwork(c *gin.Context) {
	var req struct {
		NodeID    int64  `json:"nodeId"`
		NetworkID string `json:"networkId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if req.NodeID == 0 || req.NetworkID == "" {
		response.FailWithMessage("节点ID和网络ID不能为空", c)
		return
	}

	err := networkService.RemoveNetwork(c, req.NodeID, req.NetworkID)
	if err != nil {
		global.GVA_LOG.Error("删除网络失败", zap.Error(err))
		response.FailWithMessage("删除网络失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(nil, "删除成功", c)
}
