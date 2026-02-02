import service from '@/utils/request'
// @Tags MirrorRepository
// @Summary 创建镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MirrorRepository true "创建镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /mirror/createMirrorRepository [post]
export const createMirrorRepository = (data) => {
  return service({
    url: '/mirror/createMirrorRepository',
    method: 'post',
    data
  })
}

// @Tags MirrorRepository
// @Summary 删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MirrorRepository true "删除镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mirror/deleteMirrorRepository [delete]
export const deleteMirrorRepository = (params) => {
  return service({
    url: '/mirror/deleteMirrorRepository',
    method: 'delete',
    params
  })
}

// @Tags MirrorRepository
// @Summary 批量删除镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mirror/deleteMirrorRepository [delete]
export const deleteMirrorRepositoryByIds = (params) => {
  return service({
    url: '/mirror/deleteMirrorRepositoryByIds',
    method: 'delete',
    params
  })
}

// @Tags MirrorRepository
// @Summary 更新镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.MirrorRepository true "更新镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /mirror/updateMirrorRepository [put]
export const updateMirrorRepository = (data) => {
  return service({
    url: '/mirror/updateMirrorRepository',
    method: 'put',
    data
  })
}

// @Tags MirrorRepository
// @Summary 用id查询镜像库
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.MirrorRepository true "用id查询镜像库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /mirror/findMirrorRepository [get]
export const findMirrorRepository = (params) => {
  return service({
    url: '/mirror/findMirrorRepository',
    method: 'get',
    params
  })
}

// @Tags MirrorRepository
// @Summary 分页获取镜像库列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取镜像库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mirror/getMirrorRepositoryList [get]
export const getMirrorRepositoryList = (params) => {
  return service({
    url: '/mirror/getMirrorRepositoryList',
    method: 'get',
    params
  })
}

// @Tags MirrorRepository
// @Summary 不需要鉴权的镜像库接口
// @Accept application/json
// @Produce application/json
// @Param data query cloudReq.MirrorRepositorySearch true "分页获取镜像库列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mirror/getMirrorRepositoryPublic [get]
export const getMirrorRepositoryPublic = () => {
  return service({
    url: '/mirror/getMirrorRepositoryPublic',
    method: 'get',
  })
}
