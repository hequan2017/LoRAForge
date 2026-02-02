package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type InstanceService struct{}

// CreateInstance 创建实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) CreateInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	status := "Running"
	inst.ContainerStatus = &status
	container := fmt.Sprintf("container-%d", time.Now().UnixNano())
	inst.DockerContainer = &container
	err = global.GVA_DB.Create(inst).Error
	return err
}

// CloseInstance 关闭实例
func (instService *InstanceService) CloseInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	status := "Closed"
	return global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Update("container_status", status).Error
}

// RestartInstance 重启实例
func (instService *InstanceService) RestartInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	status := "Running"
	return global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Update("container_status", status).Error
}

// DeleteInstance 删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstance(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.Instance{}, "id = ?", ID).Error
	return err
}

// DeleteInstanceByIds 批量删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstanceByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.Instance{}, "id in ?", IDs).Error
	return err
}

// UpdateInstance 更新实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) UpdateInstance(ctx context.Context, inst cloud.Instance) (err error) {
	err = global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Updates(&inst).Error
	return err
}

// GetInstance 根据ID获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstance(ctx context.Context, ID string) (inst cloud.Instance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error
	return
}

// GetInstanceInfoList 分页获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstanceInfoList(ctx context.Context, info cloudReq.InstanceSearch) (list []cloud.Instance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cloud.Instance{})
	var insts []cloud.Instance
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
    if info.InstanceName != "" {
        db = db.Where("instance_name LIKE ?", "%"+info.InstanceName+"%")
    }
    if info.MirrorId != nil {
        db = db.Where("mirror_id = ?", *info.MirrorId)
    }
    if info.TemplateId != nil {
        db = db.Where("template_id = ?", *info.TemplateId)
    }
    if info.NodeId != nil {
        db = db.Where("node_id = ?", *info.NodeId)
    }

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&insts).Error
	return insts, total, err
}
func (instService *InstanceService) GetInstanceDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	mirrorId := make([]map[string]any, 0)

	global.GVA_DB.Table("mirror_repositories").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&mirrorId)
	res["mirrorId"] = mirrorId
	nodeId := make([]map[string]any, 0)

	global.GVA_DB.Table("compute_nodes").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&nodeId)
	res["nodeId"] = nodeId
	templateId := make([]map[string]any, 0)

	global.GVA_DB.Table("product_specs").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&templateId)
	res["templateId"] = templateId
	return
}
func (instService *InstanceService) GetInstancePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
