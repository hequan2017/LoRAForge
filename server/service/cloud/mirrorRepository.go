
package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type MirrorRepositoryService struct {}
// CreateMirrorRepository 创建镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService) CreateMirrorRepository(ctx context.Context, mirror *cloud.MirrorRepository) (err error) {
	err = global.GVA_DB.Create(mirror).Error
	return err
}

// DeleteMirrorRepository 删除镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService)DeleteMirrorRepository(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.MirrorRepository{},"id = ?",ID).Error
	return err
}

// DeleteMirrorRepositoryByIds 批量删除镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService)DeleteMirrorRepositoryByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.MirrorRepository{},"id in ?",IDs).Error
	return err
}

// UpdateMirrorRepository 更新镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService)UpdateMirrorRepository(ctx context.Context, mirror cloud.MirrorRepository) (err error) {
	err = global.GVA_DB.Model(&cloud.MirrorRepository{}).Where("id = ?",mirror.ID).Updates(&mirror).Error
	return err
}

// GetMirrorRepository 根据ID获取镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService)GetMirrorRepository(ctx context.Context, ID string) (mirror cloud.MirrorRepository, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&mirror).Error
	return
}
// GetMirrorRepositoryInfoList 分页获取镜像库记录
// Author [yourname](https://github.com/yourname)
func (mirrorService *MirrorRepositoryService)GetMirrorRepositoryInfoList(ctx context.Context, info cloudReq.MirrorRepositorySearch) (list []cloud.MirrorRepository, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&cloud.MirrorRepository{})
    var mirrors []cloud.MirrorRepository
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    if info.Name != "" {
        db = db.Where("name LIKE ?", "%"+info.Name+"%")
    }
    if info.Source != "" {
        db = db.Where("source LIKE ?", "%"+info.Source+"%")
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

	err = db.Find(&mirrors).Error
	return  mirrors, total, err
}
func (mirrorService *MirrorRepositoryService)GetMirrorRepositoryPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
