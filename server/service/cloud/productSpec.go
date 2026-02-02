
package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
    cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type ProductSpecService struct {}
// CreateProductSpec 创建产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService) CreateProductSpec(ctx context.Context, spec *cloud.ProductSpec) (err error) {
	err = global.GVA_DB.Create(spec).Error
	return err
}

// DeleteProductSpec 删除产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService)DeleteProductSpec(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.ProductSpec{},"id = ?",ID).Error
	return err
}

// DeleteProductSpecByIds 批量删除产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService)DeleteProductSpecByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.ProductSpec{},"id in ?",IDs).Error
	return err
}

// UpdateProductSpec 更新产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService)UpdateProductSpec(ctx context.Context, spec cloud.ProductSpec) (err error) {
	err = global.GVA_DB.Model(&cloud.ProductSpec{}).Where("id = ?",spec.ID).Updates(&spec).Error
	return err
}

// GetProductSpec 根据ID获取产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService)GetProductSpec(ctx context.Context, ID string) (spec cloud.ProductSpec, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&spec).Error
	return
}
// GetProductSpecInfoList 分页获取产品规格记录
// Author [yourname](https://github.com/yourname)
func (specService *ProductSpecService)GetProductSpecInfoList(ctx context.Context, info cloudReq.ProductSpecSearch) (list []cloud.ProductSpec, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&cloud.ProductSpec{})
    var specs []cloud.ProductSpec
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    if info.Name != "" {
        db = db.Where("name LIKE ?", "%"+info.Name+"%")
    }
    if info.GpuModel != "" {
        db = db.Where("gpu_model LIKE ?", "%"+info.GpuModel+"%")
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

	err = db.Find(&specs).Error
	return  specs, total, err
}
func (specService *ProductSpecService)GetProductSpecPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
