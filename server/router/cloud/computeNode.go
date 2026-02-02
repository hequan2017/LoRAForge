package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ComputeNodeRouter struct {}

// InitComputeNodeRouter 初始化 算力节点 路由信息
func (s *ComputeNodeRouter) InitComputeNodeRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	nodeRouter := Router.Group("node").Use(middleware.OperationRecord())
	nodeRouterWithoutRecord := Router.Group("node")
	nodeRouterWithoutAuth := PublicRouter.Group("node")
	{
		nodeRouter.POST("createComputeNode", nodeApi.CreateComputeNode)   // 新建算力节点
		nodeRouter.DELETE("deleteComputeNode", nodeApi.DeleteComputeNode) // 删除算力节点
		nodeRouter.DELETE("deleteComputeNodeByIds", nodeApi.DeleteComputeNodeByIds) // 批量删除算力节点
		nodeRouter.PUT("updateComputeNode", nodeApi.UpdateComputeNode)    // 更新算力节点
	}
	{
		nodeRouterWithoutRecord.GET("findComputeNode", nodeApi.FindComputeNode)        // 根据ID获取算力节点
		nodeRouterWithoutRecord.GET("getComputeNodeList", nodeApi.GetComputeNodeList)  // 获取算力节点列表
	}
	{
	    nodeRouterWithoutAuth.GET("getComputeNodePublic", nodeApi.GetComputeNodePublic)  // 算力节点开放接口
	}
}
