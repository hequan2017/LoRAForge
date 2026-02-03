package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type NetworkRouter struct{}

func (r *NetworkRouter) InitNetworkRouter(Router *gin.RouterGroup) {
	networkRouter := Router.Group("cloud/network").Use(middleware.OperationRecord())
	networkRouterWithoutRecord := Router.Group("cloud/network")
	var networkApi = v1.ApiGroupApp.CloudApiGroup.NetworkApi
	{
		networkRouter.POST("create", networkApi.CreateNetwork)  // 创建网络
		networkRouter.DELETE("delete", networkApi.RemoveNetwork) // 删除网络
	}
	{
		networkRouterWithoutRecord.GET("list", networkApi.GetNetworks) // 获取网络列表
	}
}
