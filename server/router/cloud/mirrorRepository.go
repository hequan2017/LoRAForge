package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MirrorRepositoryRouter struct {}

// InitMirrorRepositoryRouter 初始化 镜像库 路由信息
func (s *MirrorRepositoryRouter) InitMirrorRepositoryRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	mirrorRouter := Router.Group("mirror").Use(middleware.OperationRecord())
	mirrorRouterWithoutRecord := Router.Group("mirror")
	mirrorRouterWithoutAuth := PublicRouter.Group("mirror")
	{
		mirrorRouter.POST("createMirrorRepository", mirrorApi.CreateMirrorRepository)   // 新建镜像库
		mirrorRouter.DELETE("deleteMirrorRepository", mirrorApi.DeleteMirrorRepository) // 删除镜像库
		mirrorRouter.DELETE("deleteMirrorRepositoryByIds", mirrorApi.DeleteMirrorRepositoryByIds) // 批量删除镜像库
		mirrorRouter.PUT("updateMirrorRepository", mirrorApi.UpdateMirrorRepository)    // 更新镜像库
	}
	{
		mirrorRouterWithoutRecord.GET("findMirrorRepository", mirrorApi.FindMirrorRepository)        // 根据ID获取镜像库
		mirrorRouterWithoutRecord.GET("getMirrorRepositoryList", mirrorApi.GetMirrorRepositoryList)  // 获取镜像库列表
	}
	{
	    mirrorRouterWithoutAuth.GET("getMirrorRepositoryPublic", mirrorApi.GetMirrorRepositoryPublic)  // 镜像库开放接口
	}
}
