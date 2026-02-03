package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SwiftWebUIRouter struct {}

// InitSwiftWebUIRouter 初始化 Swift WebUI管理 路由信息
func (s *SwiftWebUIRouter) InitSwiftWebUIRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	SwiftWebUIRouter := Router.Group("SwiftWebUI").Use(middleware.OperationRecord())
	SwiftWebUIRouterWithoutRecord := Router.Group("SwiftWebUI")
	SwiftWebUIRouterWithoutAuth := PublicRouter.Group("SwiftWebUI")
	{
		SwiftWebUIRouter.POST("createSwiftWebUI", SwiftWebUIApi.CreateSwiftWebUI)   // 新建Swift WebUI管理
		SwiftWebUIRouter.DELETE("deleteSwiftWebUI", SwiftWebUIApi.DeleteSwiftWebUI) // 删除Swift WebUI管理
		SwiftWebUIRouter.DELETE("deleteSwiftWebUIByIds", SwiftWebUIApi.DeleteSwiftWebUIByIds) // 批量删除Swift WebUI管理
		SwiftWebUIRouter.PUT("updateSwiftWebUI", SwiftWebUIApi.UpdateSwiftWebUI)    // 更新Swift WebUI管理
	}
	{
		SwiftWebUIRouterWithoutRecord.GET("findSwiftWebUI", SwiftWebUIApi.FindSwiftWebUI)        // 根据ID获取Swift WebUI管理
		SwiftWebUIRouterWithoutRecord.GET("getSwiftWebUIList", SwiftWebUIApi.GetSwiftWebUIList)  // 获取Swift WebUI管理列表
	}
	{
	    SwiftWebUIRouterWithoutAuth.GET("getSwiftWebUIDataSource", SwiftWebUIApi.GetSwiftWebUIDataSource)  // 获取Swift WebUI管理数据源
	    SwiftWebUIRouterWithoutAuth.GET("getSwiftWebUIPublic", SwiftWebUIApi.GetSwiftWebUIPublic)  // Swift WebUI管理开放接口
	}
}
