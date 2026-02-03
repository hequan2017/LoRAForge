import service from '@/utils/request'

// GetNetworks 获取网络列表
export const getNetworks = (params) => {
  return service({
    url: '/cloud/network/list',
    method: 'get',
    params
  })
}

// CreateNetwork 创建网络
export const createNetwork = (data) => {
  return service({
    url: '/cloud/network/create',
    method: 'post',
    data
  })
}

// RemoveNetwork 删除网络
export const removeNetwork = (data) => {
  return service({
    url: '/cloud/network/delete',
    method: 'delete',
    data
  })
}
