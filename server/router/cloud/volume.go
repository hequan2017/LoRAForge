package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VolumeRouter struct{}

func (r *VolumeRouter) InitVolumeRouter(Router *gin.RouterGroup) {
	volumeRouter := Router.Group("cloud/volume").Use(middleware.OperationRecord())
	volumeRouterWithoutRecord := Router.Group("cloud/volume")
	var volumeApi = v1.ApiGroupApp.CloudApiGroup.VolumeApi
	{
		volumeRouter.POST("create", volumeApi.CreateVolume)  // 创建卷
		volumeRouter.DELETE("delete", volumeApi.RemoveVolume) // 删除卷
	}
	{
		volumeRouterWithoutRecord.GET("list", volumeApi.GetVolumes) // 获取卷列表
	}
}
