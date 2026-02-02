package cloud

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type initApi struct{}

// 确保在 system 包初始化之后执行
const initOrderApi = system.InitOrderSystem + 10

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i *initApi) InitializerName() string {
	return "cloud_api"
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		// 实例管理
		{ApiGroup: "instance", Method: "POST", Path: "/inst/createInstance", Description: "新建实例管理"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/closeInstance", Description: "关闭实例"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/restartInstance", Description: "重启实例"},
		{ApiGroup: "instance", Method: "DELETE", Path: "/inst/deleteInstance", Description: "删除实例管理"},
		{ApiGroup: "instance", Method: "DELETE", Path: "/inst/deleteInstanceByIds", Description: "批量删除实例管理"},
		{ApiGroup: "instance", Method: "PUT", Path: "/inst/updateInstance", Description: "更新实例管理"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/findInstance", Description: "根据ID获取实例管理"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/getInstanceList", Description: "获取实例管理列表"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/getInstanceDataSource", Description: "获取实例管理数据源"},

		// 产品规格
		{ApiGroup: "productSpec", Method: "POST", Path: "/spec/createProductSpec", Description: "新建产品规格"},
		{ApiGroup: "productSpec", Method: "DELETE", Path: "/spec/deleteProductSpec", Description: "删除产品规格"},
		{ApiGroup: "productSpec", Method: "DELETE", Path: "/spec/deleteProductSpecByIds", Description: "批量删除产品规格"},
		{ApiGroup: "productSpec", Method: "PUT", Path: "/spec/updateProductSpec", Description: "更新产品规格"},
		{ApiGroup: "productSpec", Method: "GET", Path: "/spec/findProductSpec", Description: "根据ID获取产品规格"},
		{ApiGroup: "productSpec", Method: "GET", Path: "/spec/getProductSpecList", Description: "获取产品规格列表"},

		// 算力节点
		{ApiGroup: "computeNode", Method: "POST", Path: "/node/createComputeNode", Description: "新建算力节点"},
		{ApiGroup: "computeNode", Method: "DELETE", Path: "/node/deleteComputeNode", Description: "删除算力节点"},
		{ApiGroup: "computeNode", Method: "DELETE", Path: "/node/deleteComputeNodeByIds", Description: "批量删除算力节点"},
		{ApiGroup: "computeNode", Method: "PUT", Path: "/node/updateComputeNode", Description: "更新算力节点"},
		{ApiGroup: "computeNode", Method: "GET", Path: "/node/findComputeNode", Description: "根据ID获取算力节点"},
		{ApiGroup: "computeNode", Method: "GET", Path: "/node/getComputeNodeList", Description: "获取算力节点列表"},

		// 镜像库
		{ApiGroup: "mirrorRepository", Method: "POST", Path: "/mirror/createMirrorRepository", Description: "新建镜像库"},
		{ApiGroup: "mirrorRepository", Method: "DELETE", Path: "/mirror/deleteMirrorRepository", Description: "删除镜像库"},
		{ApiGroup: "mirrorRepository", Method: "DELETE", Path: "/mirror/deleteMirrorRepositoryByIds", Description: "批量删除镜像库"},
		{ApiGroup: "mirrorRepository", Method: "PUT", Path: "/mirror/updateMirrorRepository", Description: "更新镜像库"},
		{ApiGroup: "mirrorRepository", Method: "GET", Path: "/mirror/findMirrorRepository", Description: "根据ID获取镜像库"},
		{ApiGroup: "mirrorRepository", Method: "GET", Path: "/mirror/getMirrorRepositoryList", Description: "获取镜像库列表"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
	}
	return ctx, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/inst/createInstance", "POST").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
