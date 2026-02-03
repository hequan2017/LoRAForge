package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FineTuneTaskRouter struct{}

// InitFineTuneTaskRouter 初始化微调任务路由
func (s *FineTuneTaskRouter) InitFineTuneTaskRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	finetuneRouter := Router.Group("finetune").Use(middleware.OperationRecord())
	finetuneRouterWithoutRecord := Router.Group("finetune")
	finetuneRouterWithoutAuth := PublicRouter.Group("finetune")

	{
		// 需要操作记录的路由
		finetuneRouter.POST("createFineTuneTask", fineTuneTaskApi.CreateFineTuneTask)    // 创建微调任务
		finetuneRouter.POST("startFineTuneTask", fineTuneTaskApi.StartFineTuneTask)      // 启动微调任务
		finetuneRouter.POST("restartFineTuneTask", fineTuneTaskApi.RestartFineTuneTask)  // 重启微调任务
		finetuneRouter.POST("stopFineTuneTask", fineTuneTaskApi.StopFineTuneTask)        // 停止微调任务
		finetuneRouter.POST("cancelFineTuneTask", fineTuneTaskApi.CancelFineTuneTask)    // 取消微调任务
		finetuneRouter.DELETE("deleteFineTuneTask", fineTuneTaskApi.DeleteFineTuneTask)  // 删除微调任务
		finetuneRouter.DELETE("deleteFineTuneTaskByIds", fineTuneTaskApi.DeleteFineTuneTaskByIds) // 批量删除
		finetuneRouter.PUT("updateFineTuneTask", fineTuneTaskApi.UpdateFineTuneTask)     // 更新微调任务
		finetuneRouter.POST("createFineTuneTaskSnapshot", fineTuneTaskApi.CreateFineTuneTaskSnapshot) // 创建训练快照
		finetuneRouter.POST("exportFineTuneTaskModel", fineTuneTaskApi.ExportFineTuneTaskModel) // 导出模型
		finetuneRouter.POST("syncFineTuneTaskStatus", fineTuneTaskApi.SyncFineTuneTaskStatus) // 同步单个任务状态
		finetuneRouter.POST("syncAllFineTuneTaskStatus", fineTuneTaskApi.SyncAllFineTuneTaskStatus) // 同步所有任务状态
	}
	{
		// 不需要操作记录的路由
		finetuneRouterWithoutRecord.GET("findFineTuneTask", fineTuneTaskApi.FindFineTuneTask)        // 根据ID获取任务
		finetuneRouterWithoutRecord.GET("getFineTuneTaskList", fineTuneTaskApi.GetFineTuneTaskList) // 获取任务列表
		finetuneRouterWithoutRecord.GET("getFineTuneTaskLogs", fineTuneTaskApi.GetFineTuneTaskLogs) // 获取任务日志
		finetuneRouterWithoutRecord.GET("getFineTuneTaskMetrics", fineTuneTaskApi.GetFineTuneTaskMetrics) // 获取训练指标
		finetuneRouterWithoutRecord.GET("getFineTuneTaskSnapshots", fineTuneTaskApi.GetFineTuneTaskSnapshots) // 获取训练快照
		finetuneRouterWithoutRecord.GET("getFineTuneTaskStatistics", fineTuneTaskApi.GetFineTuneTaskStatistics) // 获取统计信息
	}
	{
		// 不需要鉴权的公开路由
		finetuneRouterWithoutAuth.GET("getFineTuneTaskDataSource", fineTuneTaskApi.GetFineTuneTaskDataSource) // 获取数据源
	}
}
