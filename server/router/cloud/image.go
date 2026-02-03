package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ImageRouter struct{}

func (r *ImageRouter) InitImageRouter(Router *gin.RouterGroup) {
	imageRouter := Router.Group("cloud/image").Use(middleware.OperationRecord())
	imageRouterWithoutRecord := Router.Group("cloud/image")
	var imageApi = v1.ApiGroupApp.CloudApiGroup.ImageApi
	{
		imageRouter.DELETE("delete", imageApi.RemoveImage) // 删除镜像
	}
	{
		imageRouterWithoutRecord.GET("list", imageApi.GetImages) // 获取镜像列表
		imageRouterWithoutRecord.GET("pull", imageApi.PullImage) // 拉取镜像
	}
}
