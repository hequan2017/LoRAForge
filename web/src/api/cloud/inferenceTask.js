import service from '@/utils/request'
// @Tags InferenceTask
// @Summary 创建AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.InferenceTask true "创建AI推理任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /inference/createInferenceTask [post]
export const createInferenceTask = (data) => {
  return service({
    url: '/inference/createInferenceTask',
    method: 'post',
    data
  })
}

// @Tags InferenceTask
// @Summary 删除AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.InferenceTask true "删除AI推理任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /inference/deleteInferenceTask [delete]
export const deleteInferenceTask = (params) => {
  return service({
    url: '/inference/deleteInferenceTask',
    method: 'delete',
    params
  })
}

// @Tags InferenceTask
// @Summary 批量删除AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除AI推理任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /inference/deleteInferenceTask [delete]
export const deleteInferenceTaskByIds = (params) => {
  return service({
    url: '/inference/deleteInferenceTaskByIds',
    method: 'delete',
    params
  })
}

// @Tags InferenceTask
// @Summary 更新AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.InferenceTask true "更新AI推理任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /inference/updateInferenceTask [put]
export const updateInferenceTask = (data) => {
  return service({
    url: '/inference/updateInferenceTask',
    method: 'put',
    data
  })
}

// @Tags InferenceTask
// @Summary 用id查询AI推理任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.InferenceTask true "用id查询AI推理任务"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /inference/findInferenceTask [get]
export const findInferenceTask = (params) => {
  return service({
    url: '/inference/findInferenceTask',
    method: 'get',
    params
  })
}

// @Tags InferenceTask
// @Summary 分页获取AI推理任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取AI推理任务列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /inference/getInferenceTaskList [get]
export const getInferenceTaskList = (params) => {
  return service({
    url: '/inference/getInferenceTaskList',
    method: 'get',
    params
  })
}

// @Tags InferenceTask
// @Summary 不需要鉴权的AI推理任务接口
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.InferenceTaskSearch true "分页获取AI推理任务列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /inference/getInferenceTaskPublic [get]
export const getInferenceTaskPublic = () => {
  return service({
    url: '/inference/getInferenceTaskPublic',
    method: 'get',
  })
}
