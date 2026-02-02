
package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type ComputeNodeService struct {}
// CreateComputeNode 创建算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService) CreateComputeNode(ctx context.Context, node *cloud.ComputeNode) (err error) {
	err = global.GVA_DB.Create(node).Error
	return err
}

// DeleteComputeNode 删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService)DeleteComputeNode(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.ComputeNode{},"id = ?",ID).Error
	return err
}

// DeleteComputeNodeByIds 批量删除算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService)DeleteComputeNodeByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.ComputeNode{},"id in ?",IDs).Error
	return err
}

// UpdateComputeNode 更新算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService)UpdateComputeNode(ctx context.Context, node cloud.ComputeNode) (err error) {
	err = global.GVA_DB.Model(&cloud.ComputeNode{}).Where("id = ?",node.ID).Updates(&node).Error
	return err
}

// GetComputeNode 根据ID获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService)GetComputeNode(ctx context.Context, ID string) (node cloud.ComputeNode, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&node).Error
	return
}
// GetComputeNodeInfoList 分页获取算力节点记录
// Author [yourname](https://github.com/yourname)
func (nodeService *ComputeNodeService)GetComputeNodeInfoList(ctx context.Context, info cloudReq.ComputeNodeSearch) (list []cloud.ComputeNode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&cloud.ComputeNode{})
    var nodes []cloud.ComputeNode
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    if info.Name != "" {
        db = db.Where("name LIKE ?", "%"+info.Name+"%")
    }
    if info.Region != "" {
        db = db.Where("region LIKE ?", "%"+info.Region+"%")
    }
    if info.PublicIp != "" {
        db = db.Where("public_ip LIKE ?", "%"+info.PublicIp+"%")
    }
    if info.GpuName != "" {
        db = db.Where("gpu_name LIKE ?", "%"+info.GpuName+"%")
    }
    if info.IsListed != nil {
        db = db.Where("is_listed = ?", *info.IsListed)
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&nodes).Error
	return  nodes, total, err
}
func (nodeService *ComputeNodeService)GetComputeNodePublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
