import service from '@/utils/request'

// GetVolumes 获取卷列表
export const getVolumes = (params) => {
  return service({
    url: '/cloud/volume/list',
    method: 'get',
    params
  })
}

// CreateVolume 创建卷
export const createVolume = (data) => {
  return service({
    url: '/cloud/volume/create',
    method: 'post',
    data
  })
}

// RemoveVolume 删除卷
export const removeVolume = (data) => {
  return service({
    url: '/cloud/volume/delete',
    method: 'delete',
    data
  })
}
