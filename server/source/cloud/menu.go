package cloud

import (
	"context"

	. "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = system.InitOrderSystem + 11

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i *initMenu) InitializerName() string {
	return "cloud_menu"
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("name = ?", "cloud").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 1. 创建父级菜单
	parentMenu := SysBaseMenu{
		MenuLevel: 0,
		Hidden:    false,
		ParentId:  0,
		Path:      "cloud",
		Name:      "cloud",
		Component: "view/routerHolder.vue",
		Sort:      10,
		Meta:      Meta{Title: "云资源管理", Icon: "cloudy"},
	}

	// 检查是否存在，如果存在则获取ID，不存在则创建
	var existingParent SysBaseMenu
	if err := db.Where("name = ?", parentMenu.Name).First(&existingParent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Create(&parentMenu).Error; err != nil {
				return ctx, errors.Wrap(err, "创建云资源管理父菜单失败")
			}
			existingParent = parentMenu
		} else {
			return ctx, err
		}
	}

	// 2. 创建子菜单
	childMenus := []SysBaseMenu{
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "instance",
			Name:      "instance",
			Component: "view/cloud/instance/instance.vue",
			Sort:      1,
			Meta:      Meta{Title: "实例管理", Icon: "monitor"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "computeNode",
			Name:      "computeNode",
			Component: "view/cloud/compute_node/compute_node.vue",
			Sort:      2,
			Meta:      Meta{Title: "算力节点", Icon: "cpu"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "productSpec",
			Name:      "productSpec",
			Component: "view/cloud/product_spec/product_spec.vue",
			Sort:      3,
			Meta:      Meta{Title: "产品规格", Icon: "files"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "nodeImage",
			Name:      "nodeImage",
			Component: "view/cloud/image/image.vue",
			Sort:      4,
			Meta:      Meta{Title: "节点镜像", Icon: "picture"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "nodeNetwork",
			Name:      "nodeNetwork",
			Component: "view/cloud/network/network.vue",
			Sort:      5,
			Meta:      Meta{Title: "节点网络", Icon: "connection"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "nodeVolume",
			Name:      "nodeVolume",
			Component: "view/cloud/volume/volume.vue",
			Sort:      6,
			Meta:      Meta{Title: "节点数据卷", Icon: "files"},
		},
		{
			MenuLevel: 1,
			Hidden:    false,
			ParentId:  existingParent.ID,
			Path:      "finetuneTask",
			Name:      "finetuneTask",
			Component: "view/cloud/finetune_task/finetune_task.vue",
			Sort:      7,
			Meta:      Meta{Title: "微调任务", Icon: "aim"},
		},
	}

	for _, menu := range childMenus {
		var existingMenu SysBaseMenu
		if err := db.Where("name = ?", menu.Name).First(&existingMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Create(&menu).Error; err != nil {
					return ctx, errors.Wrap(err, "创建子菜单 "+menu.Meta.Title+" 失败")
				}
			} else {
				return ctx, err
			}
		} else {
			// 如果已存在，更新 Hidden 状态
			if menu.Hidden != existingMenu.Hidden {
				if err = db.Model(&existingMenu).Update("hidden", menu.Hidden).Error; err != nil {
					return ctx, errors.Wrap(err, "更新子菜单 "+menu.Meta.Title+" 失败")
				}
			}
		}
	}
	
	// 删除旧的 mirrorRepository 菜单
	if err := db.Where("name = ?", "mirrorRepository").Delete(&SysBaseMenu{}).Error; err != nil {
		// 忽略删除错误，可能已经不存在
	}

	return ctx, nil
}
