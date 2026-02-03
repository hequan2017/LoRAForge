import service from '@/utils/request'

// GetImages 获取镜像列表
// @Tags Image
// @Summary 获取镜像列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param nodeId query int true "节点ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cloud/image/list [get]
export const getImages = (params) => {
  return service({
    url: '/cloud/image/list',
    method: 'get',
    params
  })
}

// RemoveImage 删除镜像
// @Tags Image
// @Summary 删除镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body object true "删除参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cloud/image/delete [delete]
export const removeImage = (data) => {
  return service({
    url: '/cloud/image/delete',
    method: 'delete',
    data
  })
}

// PullImage 拉取镜像 (SSE)
// 注意：SSE 通常使用 EventSource 或 fetch 直接处理，这里仅定义 URL 辅助方法
export const getPullImageUrl = (nodeId, imageName) => {
    const baseUrl = import.meta.env.VITE_BASE_API || ''
    return `${baseUrl}/cloud/image/pull?nodeId=${nodeId}&imageName=${encodeURIComponent(imageName)}`
}
