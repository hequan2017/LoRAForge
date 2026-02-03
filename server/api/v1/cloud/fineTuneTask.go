package cloud

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FineTuneTaskApi struct{}

// CreateFineTuneTask 创建微调任务
// @Tags FineTuneTask
// @Summary 创建微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloud.FineTuneTask true "微调任务信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /finetune/createFineTuneTask [post]
func (api *FineTuneTaskApi) CreateFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var task cloud.FineTuneTask
	if err := c.ShouldBindJSON(&task); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置用户ID
	userID := utils.GetUserID(c)
	task.UserID = new(int64)
	*task.UserID = int64(userID)

	if err := fineTuneTaskService.CreateFineTuneTask(ctx, &task); err != nil {
		global.GVA_LOG.Error("创建微调任务失败", zap.Error(err))
		response.FailWithMessage("创建微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("创建微调任务成功", c)
}

// StartFineTuneTask 启动微调任务
// @Tags FineTuneTask
// @Summary 启动微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskStart true "启动任务请求"
// @Success 200 {object} response.Response{msg=string} "启动成功"
// @Router /finetune/startFineTuneTask [post]
func (api *FineTuneTaskApi) StartFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskStart
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := fineTuneTaskService.StartFineTuneTask(ctx, *req.ID); err != nil {
		global.GVA_LOG.Error("启动微调任务失败", zap.Error(err))
		response.FailWithMessage("启动微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("微调任务启动成功", c)
}

// StopFineTuneTask 停止微调任务
// @Tags FineTuneTask
// @Summary 停止微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskStop true "停止任务请求"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /finetune/stopFineTuneTask [post]
func (api *FineTuneTaskApi) StopFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskStop
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	if err := fineTuneTaskService.StopFineTuneTask(ctx, *req.ID, reason); err != nil {
		global.GVA_LOG.Error("停止微调任务失败", zap.Error(err))
		response.FailWithMessage("停止微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("微调任务已停止", c)
}

// RestartFineTuneTask 重启微调任务
// @Tags FineTuneTask
// @Summary 重启微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskStart true "重启任务请求"
// @Success 200 {object} response.Response{msg=string} "重启成功"
// @Router /finetune/restartFineTuneTask [post]
func (api *FineTuneTaskApi) RestartFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskStart
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := fineTuneTaskService.RestartFineTuneTask(ctx, *req.ID); err != nil {
		global.GVA_LOG.Error("重启微调任务失败", zap.Error(err))
		response.FailWithMessage("重启微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("微调任务重启成功", c)
}

// CancelFineTuneTask 取消微调任务
// @Tags FineTuneTask
// @Summary 取消微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskCancelRequest true "取消任务请求"
// @Success 200 {object} response.Response{msg=string} "取消成功"
// @Router /finetune/cancelFineTuneTask [post]
func (api *FineTuneTaskApi) CancelFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	reason := "用户取消"
	if req.Reason != nil {
		reason = *req.Reason
	}

	if err := fineTuneTaskService.StopFineTuneTask(ctx, *req.ID, reason); err != nil {
		global.GVA_LOG.Error("取消微调任务失败", zap.Error(err))
		response.FailWithMessage("取消微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("微调任务已取消", c)
}

// DeleteFineTuneTask 删除微调任务
// @Tags FineTuneTask
// @Summary 删除微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskDelete true "删除任务请求"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /finetune/deleteFineTuneTask [delete]
func (api *FineTuneTaskApi) DeleteFineTuneTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := fineTuneTaskService.DeleteFineTuneTask(ctx, *req.ID); err != nil {
		global.GVA_LOG.Error("删除微调任务失败", zap.Error(err))
		response.FailWithMessage("删除微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除微调任务成功", c)
}

// DeleteFineTuneTaskByIds 批量删除微调任务
// @Tags FineTuneTask
// @Summary 批量删除微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskBatchDelete true "批量删除请求"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /finetune/deleteFineTuneTaskByIds [delete]
func (api *FineTuneTaskApi) DeleteFineTuneTaskByIds(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskBatchDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := fineTuneTaskService.DeleteFineTuneTaskByIds(ctx, req.IDs); err != nil {
		global.GVA_LOG.Error("批量删除微调任务失败", zap.Error(err))
		response.FailWithMessage("批量删除微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("批量删除微调任务成功", c)
}

// UpdateFineTuneTask 更新微调任务
// @Tags FineTuneTask
// @Summary 更新微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskUpdate true "更新任务请求"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /finetune/updateFineTuneTask [put]
func (api *FineTuneTaskApi) UpdateFineTuneTask(c *gin.Context) {
	var req cloudReq.FineTuneTaskUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, *req.ID).Error; err != nil {
		response.FailWithMessage("任务不存在", c)
		return
	}

	// 更新允许修改的字段
	updates := make(map[string]interface{})
	if req.TaskName != nil {
		updates["task_name"] = *req.TaskName
	}
	if req.TaskDescription != nil {
		updates["task_description"] = *req.TaskDescription
	}
	if req.TaskStatus != nil {
		updates["task_status"] = *req.TaskStatus
	}
	if req.Progress != nil {
		updates["progress"] = *req.Progress
	}
	if req.CurrentStep != nil {
		updates["current_step"] = *req.CurrentStep
	}
	if req.TotalSteps != nil {
		updates["total_steps"] = *req.TotalSteps
	}
	if req.ErrorMessage != nil {
		updates["error_message"] = *req.ErrorMessage
	}
	if req.OutputModelPath != nil {
		updates["output_model_path"] = *req.OutputModelPath
	}
	if req.ValidationLoss != nil {
		updates["validation_loss"] = *req.ValidationLoss
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if err := global.GVA_DB.Model(&task).Updates(updates).Error; err != nil {
		global.GVA_LOG.Error("更新微调任务失败", zap.Error(err))
		response.FailWithMessage("更新微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新微调任务成功", c)
}

// FindFineTuneTask 根据ID获取微调任务
// @Tags FineTuneTask
// @Summary 根据ID获取微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {object} response.Response{data=cloud.FineTuneTask,msg=string} "查询成功"
// @Router /finetune/findFineTuneTask [get]
func (api *FineTuneTaskApi) FindFineTuneTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("任务ID不能为空", c)
		return
	}

	// 查询任务
	var task cloud.FineTuneTask
	if err := global.GVA_DB.Where("id = ?", id).First(&task).Error; err != nil {
		global.GVA_LOG.Error("查询微调任务失败", zap.Error(err))
		response.FailWithMessage("查询微调任务失败: "+err.Error(), c)
		return
	}

	response.OkWithData(task, c)
}

// GetFineTuneTaskList 分页获取微调任务列表
// @Tags FineTuneTask
// @Summary 分页获取微调任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.FineTuneTaskSearch true "搜索条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskList [get]
func (api *FineTuneTaskApi) GetFineTuneTaskList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo cloudReq.FineTuneTaskSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := fineTuneTaskService.GetFineTuneTaskList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取微调任务列表失败", zap.Error(err))
		response.FailWithMessage("获取微调任务列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetFineTuneTaskDataSource 获取微调任务数据源
// @Tags FineTuneTask
// @Summary 获取微调任务数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string][]map[string]any,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskDataSource [get]
func (api *FineTuneTaskApi) GetFineTuneTaskDataSource(c *gin.Context) {
	ctx := c.Request.Context()

	dataSource, err := fineTuneTaskService.GetFineTuneTaskDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取数据源失败", zap.Error(err))
		response.FailWithMessage("获取数据源失败: "+err.Error(), c)
		return
	}

	response.OkWithData(dataSource, c)
}

// GetFineTuneTaskLogs 获取微调任务日志
// @Tags FineTuneTask
// @Summary 获取微调任务日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Param offset query int false "偏移量" default(0)
// @Param limit query int false "行数限制" default(100)
// @Success 200 {object} response.Response{data=[]string,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskLogs [get]
func (api *FineTuneTaskApi) GetFineTuneTaskLogs(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("任务ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("任务ID格式错误", c)
		return
	}

	offset := 0
	limit := 100

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	logs, err := fineTuneTaskService.GetTaskLogs(c.Request.Context(), uint(id), offset, limit)
	if err != nil {
		global.GVA_LOG.Error("获取任务日志失败", zap.Error(err))
		response.FailWithMessage("获取任务日志失败: "+err.Error(), c)
		return
	}

	response.OkWithData(logs, c)
}

// GetFineTuneTaskMetrics 获取微调任务训练指标
// @Tags FineTuneTask
// @Summary 获取微调任务训练指标
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {object} response.Response{data=map[string]any,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskMetrics [get]
func (api *FineTuneTaskApi) GetFineTuneTaskMetrics(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("任务ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("任务ID格式错误", c)
		return
	}

	metrics, err := fineTuneTaskService.GetMetrics(c.Request.Context(), uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取任务指标失败", zap.Error(err))
		response.FailWithMessage("获取任务指标失败: "+err.Error(), c)
		return
	}

	response.OkWithData(metrics, c)
}

// ExportFineTuneTaskModel 导出微调模型
// @Tags FineTuneTask
// @Summary 导出微调模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskExportRequest true "导出请求"
// @Success 200 {object} response.Response{msg=string} "导出成功"
// @Router /finetune/exportFineTuneTaskModel [post]
func (api *FineTuneTaskApi) ExportFineTuneTaskModel(c *gin.Context) {
	ctx := c.Request.Context()

	var req cloudReq.FineTuneTaskExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	exportFormat := "merge"
	if req.ExportFormat != nil && *req.ExportFormat != "" {
		exportFormat = *req.ExportFormat
	}

	exportPath := ""
	if req.ExportPath != nil {
		exportPath = *req.ExportPath
	}

	quantizeBit := 4
	if req.QuantizeBit != nil {
		quantizeBit = *req.QuantizeBit
	}

	if err := fineTuneTaskService.ExportModel(ctx, *req.ID, exportFormat, exportPath, quantizeBit); err != nil {
		global.GVA_LOG.Error("导出模型失败", zap.Error(err))
		response.FailWithMessage("导出模型失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("模型导出任务已启动", c)
}

// GetFineTuneTaskStatistics 获取微调任务统计信息
// @Tags FineTuneTask
// @Summary 获取微调任务统计信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]any,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskStatistics [get]
func (api *FineTuneTaskApi) GetFineTuneTaskStatistics(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := fineTuneTaskService.GetTaskStatistics(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取统计信息失败", zap.Error(err))
		response.FailWithMessage("获取统计信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(stats, c)
}

// CreateFineTuneTaskSnapshot 创建训练快照
// @Tags FineTuneTask
// @Summary 创建训练快照
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cloudReq.FineTuneTaskSnapshotCreate true "快照数据"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /finetune/createFineTuneTaskSnapshot [post]
func (api *FineTuneTaskApi) CreateFineTuneTaskSnapshot(c *gin.Context) {
	ctx := c.Request.Context()

	var snapshot cloud.FineTuneTaskSnapshot
	if err := c.ShouldBindJSON(&snapshot); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := fineTuneTaskService.CreateSnapshot(ctx, &snapshot); err != nil {
		global.GVA_LOG.Error("创建快照失败", zap.Error(err))
		response.FailWithMessage("创建快照失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("创建快照成功", c)
}

// GetFineTuneTaskSnapshots 获取训练快照列表
// @Tags FineTuneTask
// @Summary 获取训练快照列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param taskId query uint true "任务ID"
// @Success 200 {object} response.Response{data=[]cloud.FineTuneTaskSnapshot,msg=string} "获取成功"
// @Router /finetune/getFineTuneTaskSnapshots [get]
func (api *FineTuneTaskApi) GetFineTuneTaskSnapshots(c *gin.Context) {
	taskIdStr := c.Query("taskId")
	if taskIdStr == "" {
		response.FailWithMessage("任务ID不能为空", c)
		return
	}

	taskId, err := strconv.ParseUint(taskIdStr, 10, 32)
	if err != nil {
		response.FailWithMessage("任务ID格式错误", c)
		return
	}

	snapshots, err := fineTuneTaskService.GetSnapshots(c.Request.Context(), uint(taskId))
	if err != nil {
		global.GVA_LOG.Error("获取快照列表失败", zap.Error(err))
		response.FailWithMessage("获取快照列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(snapshots, c)
}

// SyncFineTuneTaskStatus 同步微调任务状态
// @Tags FineTuneTask
// @Summary 同步微调任务状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {object} response.Response{msg=string} "同步成功"
// @Router /finetune/syncFineTuneTaskStatus [post]
func (api *FineTuneTaskApi) SyncFineTuneTaskStatus(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("任务ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("任务ID格式错误", c)
		return
	}

	if err := fineTuneTaskService.SyncSwiftTask(c.Request.Context(), uint(id)); err != nil {
		global.GVA_LOG.Error("同步任务状态失败", zap.Error(err))
		response.FailWithMessage("同步任务状态失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("同步任务状态成功", c)
}

// SyncAllFineTuneTaskStatus 同步所有微调任务状态
// @Tags FineTuneTask
// @Summary 同步所有微调任务状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "同步成功"
// @Router /finetune/syncAllFineTuneTaskStatus [post]
func (api *FineTuneTaskApi) SyncAllFineTuneTaskStatus(c *gin.Context) {
	ctx := c.Request.Context()

	if err := fineTuneTaskService.UpdateTaskStatus(ctx); err != nil {
		global.GVA_LOG.Error("同步所有任务状态失败", zap.Error(err))
		response.FailWithMessage("同步所有任务状态失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("同步所有任务状态成功", c)
}
