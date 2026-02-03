package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InferenceTaskRouter struct {}

// InitInferenceTaskRouter 初始化 AI推理任务 路由信息
func (s *InferenceTaskRouter) InitInferenceTaskRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	inferenceRouter := Router.Group("inference").Use(middleware.OperationRecord())
	inferenceRouterWithoutRecord := Router.Group("inference")
	inferenceRouterWithoutAuth := PublicRouter.Group("inference")
	{
		inferenceRouter.POST("createInferenceTask", inferenceApi.CreateInferenceTask)   // 新建AI推理任务
		inferenceRouter.DELETE("deleteInferenceTask", inferenceApi.DeleteInferenceTask) // 删除AI推理任务
		inferenceRouter.DELETE("deleteInferenceTaskByIds", inferenceApi.DeleteInferenceTaskByIds) // 批量删除AI推理任务
		inferenceRouter.PUT("updateInferenceTask", inferenceApi.UpdateInferenceTask)    // 更新AI推理任务
	}
	{
		inferenceRouterWithoutRecord.GET("findInferenceTask", inferenceApi.FindInferenceTask)        // 根据ID获取AI推理任务
		inferenceRouterWithoutRecord.GET("getInferenceTaskList", inferenceApi.GetInferenceTaskList)  // 获取AI推理任务列表
	}
	{
	    inferenceRouterWithoutAuth.GET("getInferenceTaskPublic", inferenceApi.GetInferenceTaskPublic)  // AI推理任务开放接口
	}
}
