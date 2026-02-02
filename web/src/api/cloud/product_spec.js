import service from '@/utils/request'
// @Tags ProductSpec
// @Summary 创建产品规格
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.ProductSpec true "创建产品规格"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /spec/createProductSpec [post]
export const createProductSpec = (data) => {
  return service({
    url: '/spec/createProductSpec',
    method: 'post',
    data
  })
}

// @Tags ProductSpec
// @Summary 删除产品规格
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.ProductSpec true "删除产品规格"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /spec/deleteProductSpec [delete]
export const deleteProductSpec = (params) => {
  return service({
    url: '/spec/deleteProductSpec',
    method: 'delete',
    params
  })
}

// @Tags ProductSpec
// @Summary 批量删除产品规格
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除产品规格"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /spec/deleteProductSpec [delete]
export const deleteProductSpecByIds = (params) => {
  return service({
    url: '/spec/deleteProductSpecByIds',
    method: 'delete',
    params
  })
}

// @Tags ProductSpec
// @Summary 更新产品规格
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.ProductSpec true "更新产品规格"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /spec/updateProductSpec [put]
export const updateProductSpec = (data) => {
  return service({
    url: '/spec/updateProductSpec',
    method: 'put',
    data
  })
}

// @Tags ProductSpec
// @Summary 用id查询产品规格
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.ProductSpec true "用id查询产品规格"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /spec/findProductSpec [get]
export const findProductSpec = (params) => {
  return service({
    url: '/spec/findProductSpec',
    method: 'get',
    params
  })
}

// @Tags ProductSpec
// @Summary 分页获取产品规格列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取产品规格列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /spec/getProductSpecList [get]
export const getProductSpecList = (params) => {
  return service({
    url: '/spec/getProductSpecList',
    method: 'get',
    params
  })
}

// @Tags ProductSpec
// @Summary 不需要鉴权的产品规格接口
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.ProductSpecSearch true "分页获取产品规格列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /spec/getProductSpecPublic [get]
export const getProductSpecPublic = () => {
  return service({
    url: '/spec/getProductSpecPublic',
    method: 'get',
  })
}
