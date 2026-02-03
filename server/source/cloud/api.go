package cloud

import (
	"context"

	. "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApi = system.InitOrderSystem + 12

type initApi struct{}

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
	return db.Migrator().HasTable(&SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []SysApi{
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
		
		// 实例高级功能 (新增)
		{ApiGroup: "instance", Method: "GET", Path: "/inst/stats", Description: "获取容器监控数据"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/webssh", Description: "WebSSH连接"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/logs", Description: "容器日志"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/file/list", Description: "获取容器文件列表"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/file/download", Description: "下载容器文件"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/file/upload", Description: "上传容器文件"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/file/delete", Description: "删除容器文件"},

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
		
		// 镜像管理 (新增)
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/image/list", Description: "获取节点镜像列表"},
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/image/pull", Description: "拉取镜像"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/image/delete", Description: "删除镜像"},

		// 网络管理 (新增)
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/network/list", Description: "获取网络列表"},
		{ApiGroup: "cloud", Method: "POST", Path: "/cloud/network/create", Description: "创建网络"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/network/delete", Description: "删除网络"},

		// 卷管理 (新增)
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/volume/list", Description: "获取卷列表"},
		{ApiGroup: "cloud", Method: "POST", Path: "/cloud/volume/create", Description: "创建卷"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/volume/delete", Description: "删除卷"},

		// 微调任务 (新增)
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/createFineTuneTask", Description: "创建微调任务"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/startFineTuneTask", Description: "启动微调任务"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/restartFineTuneTask", Description: "重启微调任务"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/stopFineTuneTask", Description: "停止微调任务"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/cancelFineTuneTask", Description: "取消微调任务"},
		{ApiGroup: "fineTuneTask", Method: "DELETE", Path: "/finetune/deleteFineTuneTask", Description: "删除微调任务"},
		{ApiGroup: "fineTuneTask", Method: "DELETE", Path: "/finetune/deleteFineTuneTaskByIds", Description: "批量删除微调任务"},
		{ApiGroup: "fineTuneTask", Method: "PUT", Path: "/finetune/updateFineTuneTask", Description: "更新微调任务"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/createFineTuneTaskSnapshot", Description: "创建训练快照"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/exportFineTuneTaskModel", Description: "导出模型"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/syncFineTuneTaskStatus", Description: "同步任务状态"},
		{ApiGroup: "fineTuneTask", Method: "POST", Path: "/finetune/syncAllFineTuneTaskStatus", Description: "同步所有任务状态"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/findFineTuneTask", Description: "根据ID获取任务"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskList", Description: "获取任务列表"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskLogs", Description: "获取任务日志"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskMetrics", Description: "获取训练指标"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskSnapshots", Description: "获取训练快照"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskStatistics", Description: "获取统计信息"},
		{ApiGroup: "fineTuneTask", Method: "GET", Path: "/finetune/getFineTuneTaskDataSource", Description: "获取任务数据源"},
	}

	if err := db.Create(&entities).Error; err != nil {
		// 忽略唯一键冲突错误，因为可能部分API已存在
		// 实际上 GORM Create 遇到冲突会报错，更好的做法是使用 Clauses(clause.OnConflict{DoNothing: true})
		// 但为了保持代码简单且假设初始化主要在空库或增量执行，这里仅记录日志或忽略重复
		// 更好的方式是使用 FirstOrCreate 或 循环 Check 
		// 这里简单处理：如果报错，尝试逐个插入
		for _, entity := range entities {
			db.FirstOrCreate(&entity, SysApi{Path: entity.Path, Method: entity.Method})
		}
	}
	
	// 删除旧的 mirrorRepository API
	db.Where("api_group = ?", "mirrorRepository").Delete(&SysApi{})

	return ctx, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/inst/createInstance", "POST").
		First(&SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
