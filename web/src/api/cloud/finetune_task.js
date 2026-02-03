import service from '@/utils/request'

// @Tags FineTuneTask
// @Summary 创建微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.FineTuneTask true "微调任务信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /finetune/createFineTuneTask [post]
export const createFineTuneTask = (data) => {
  return service({
    url: '/finetune/createFineTuneTask',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 启动微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskStart true "启动任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"启动成功"}"
// @Router /finetune/startFineTuneTask [post]
export const startFineTuneTask = (data) => {
  return service({
    url: '/finetune/startFineTuneTask',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 停止微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskStop true "停止任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"停止成功"}"
// @Router /finetune/stopFineTuneTask [post]
export const stopFineTuneTask = (data) => {
  return service({
    url: '/finetune/stopFineTuneTask',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 重启微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskStart true "重启任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"重启成功"}"
// @Router /finetune/restartFineTuneTask [post]
export const restartFineTuneTask = (data) => {
  return service({
    url: '/finetune/restartFineTuneTask',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 取消微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskCancelRequest true "取消任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"取消成功"}"
// @Router /finetune/cancelFineTuneTask [post]
export const cancelFineTuneTask = (data) => {
  return service({
    url: '/finetune/cancelFineTuneTask',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 删除微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskDelete true "删除任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /finetune/deleteFineTuneTask [delete]
export const deleteFineTuneTask = (data) => {
  return service({
    url: '/finetune/deleteFineTuneTask',
    method: 'delete',
    data
  })
}

// @Tags FineTuneTask
// @Summary 批量删除微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskBatchDelete true "批量删除请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /finetune/deleteFineTuneTaskByIds [delete]
export const deleteFineTuneTaskByIds = (data) => {
  return service({
    url: '/finetune/deleteFineTuneTaskByIds',
    method: 'delete',
    data
  })
}

// @Tags FineTuneTask
// @Summary 更新微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskUpdate true "更新任务请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /finetune/updateFineTuneTask [put]
export const updateFineTuneTask = (data) => {
  return service({
    url: '/finetune/updateFineTuneTask',
    method: 'put',
    data
  })
}

// @Tags FineTuneTask
// @Summary 根据ID获取微调任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /finetune/findFineTuneTask [get]
export const findFineTuneTask = (params) => {
  return service({
    url: '/finetune/findFineTuneTask',
    method: 'get',
    params
  })
}

// @Tags FineTuneTask
// @Summary 分页获取微调任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.FineTuneTaskSearch true "搜索条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskList [get]
export const getFineTuneTaskList = (params) => {
  return service({
    url: '/finetune/getFineTuneTaskList',
    method: 'get',
    params
  })
}

// @Tags FineTuneTask
// @Summary 获取微调任务数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskDataSource [get]
export const getFineTuneTaskDataSource = () => {
  return service({
    url: '/finetune/getFineTuneTaskDataSource',
    method: 'get'
  })
}

// @Tags FineTuneTask
// @Summary 获取微调任务日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Param offset query int false "偏移量"
// @Param limit query int false "行数限制"
// @Success 200 {string} string "{"success":true,"data":[],"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskLogs [get]
export const getFineTuneTaskLogs = (params) => {
  return service({
    url: '/finetune/getFineTuneTaskLogs',
    method: 'get',
    params
  })
}

// @Tags FineTuneTask
// @Summary 获取微调任务训练指标
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskMetrics [get]
export const getFineTuneTaskMetrics = (params) => {
  return service({
    url: '/finetune/getFineTuneTaskMetrics',
    method: 'get',
    params
  })
}

// @Tags FineTuneTask
// @Summary 导出微调模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.FineTuneTaskExportRequest true "导出请求"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导出成功"}"
// @Router /finetune/exportFineTuneTaskModel [post]
export const exportFineTuneTaskModel = (data) => {
  return service({
    url: '/finetune/exportFineTuneTaskModel',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 获取微调任务统计信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskStatistics [get]
export const getFineTuneTaskStatistics = () => {
  return service({
    url: '/finetune/getFineTuneTaskStatistics',
    method: 'get'
  })
}

// @Tags FineTuneTask
// @Summary 创建训练快照
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.FineTuneTaskSnapshot true "快照数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /finetune/createFineTuneTaskSnapshot [post]
export const createFineTuneTaskSnapshot = (data) => {
  return service({
    url: '/finetune/createFineTuneTaskSnapshot',
    method: 'post',
    data
  })
}

// @Tags FineTuneTask
// @Summary 获取训练快照列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param taskId query uint true "任务ID"
// @Success 200 {string} string "{"success":true,"data":[],"msg":"获取成功"}"
// @Router /finetune/getFineTuneTaskSnapshots [get]
export const getFineTuneTaskSnapshots = (params) => {
  return service({
    url: '/finetune/getFineTuneTaskSnapshots',
    method: 'get',
    params
  })
}

// @Tags FineTuneTask
// @Summary 同步微调任务状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "任务ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"同步成功"}"
// @Router /finetune/syncFineTuneTaskStatus [post]
export const syncFineTuneTaskStatus = (params) => {
  return service({
    url: '/finetune/syncFineTuneTaskStatus',
    method: 'post',
    params
  })
}

// @Tags FineTuneTask
// @Summary 同步所有微调任务状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"同步成功"}"
// @Router /finetune/syncAllFineTuneTaskStatus [post]
export const syncAllFineTuneTaskStatus = () => {
  return service({
    url: '/finetune/syncAllFineTuneTaskStatus',
    method: 'post'
  })
}
