import service from '@/utils/request'
// @Tags SwiftWebUI
// @Summary 创建Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SwiftWebUI true "创建Swift WebUI管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /SwiftWebUI/createSwiftWebUI [post]
export const createSwiftWebUI = (data) => {
  return service({
    url: '/SwiftWebUI/createSwiftWebUI',
    method: 'post',
    data
  })
}

// @Tags SwiftWebUI
// @Summary 删除Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SwiftWebUI true "删除Swift WebUI管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SwiftWebUI/deleteSwiftWebUI [delete]
export const deleteSwiftWebUI = (params) => {
  return service({
    url: '/SwiftWebUI/deleteSwiftWebUI',
    method: 'delete',
    params
  })
}

// @Tags SwiftWebUI
// @Summary 批量删除Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Swift WebUI管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SwiftWebUI/deleteSwiftWebUI [delete]
export const deleteSwiftWebUIByIds = (params) => {
  return service({
    url: '/SwiftWebUI/deleteSwiftWebUIByIds',
    method: 'delete',
    params
  })
}

// @Tags SwiftWebUI
// @Summary 更新Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SwiftWebUI true "更新Swift WebUI管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /SwiftWebUI/updateSwiftWebUI [put]
export const updateSwiftWebUI = (data) => {
  return service({
    url: '/SwiftWebUI/updateSwiftWebUI',
    method: 'put',
    data
  })
}

// @Tags SwiftWebUI
// @Summary 用id查询Swift WebUI管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SwiftWebUI true "用id查询Swift WebUI管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SwiftWebUI/findSwiftWebUI [get]
export const findSwiftWebUI = (params) => {
  return service({
    url: '/SwiftWebUI/findSwiftWebUI',
    method: 'get',
    params
  })
}

// @Tags SwiftWebUI
// @Summary 分页获取Swift WebUI管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Swift WebUI管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /SwiftWebUI/getSwiftWebUIList [get]
export const getSwiftWebUIList = (params) => {
  return service({
    url: '/SwiftWebUI/getSwiftWebUIList',
    method: 'get',
    params
  })
}
// @Tags SwiftWebUI
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SwiftWebUI/findSwiftWebUIDataSource [get]
export const getSwiftWebUIDataSource = () => {
  return service({
    url: '/SwiftWebUI/getSwiftWebUIDataSource',
    method: 'get',
  })
}

// @Tags SwiftWebUI
// @Summary 不需要鉴权的Swift WebUI管理接口
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.SwiftWebUISearch true "分页获取Swift WebUI管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /SwiftWebUI/getSwiftWebUIPublic [get]
export const getSwiftWebUIPublic = () => {
  return service({
    url: '/SwiftWebUI/getSwiftWebUIPublic',
    method: 'get',
  })
}
