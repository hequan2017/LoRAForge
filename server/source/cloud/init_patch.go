package cloud

import (
	"context"

	adapter "github.com/casbin/gorm-adapter/v3"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// 设置一个较高的执行顺序，确保在基础表创建之后执行
const initOrderPatch = system.InitOrderSystem + 99

type initPatch struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPatch, &initPatch{})
}

func (i *initPatch) InitializerName() string {
	// 使用一个新的名称，确保未被执行过
	return "cloud_permission_patch_v1"
}

func (i *initPatch) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initPatch) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	// 确保 api 表和 casbin_rule 表都存在
	return db.Migrator().HasTable(&sysModel.SysApi{}) && db.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i *initPatch) DataInserted(ctx context.Context) bool {
	// 始终返回 false，让 InitializeData 执行检查和插入逻辑
	// 因为我们是补丁，内部会自己判断是否需要插入
	// 但为了避免重复执行 patch 本身（GVA 框架会记录 InitializerName），这里返回 false 是安全的，
	// 只要 InitializerName 是唯一的，GVA 就会调用 InitializeData 一次。
	// GVA 的逻辑是：检查 sys_init 表中是否有 InitializerName，如果有，跳过；
	// 如果没有，调用 DataInserted；如果 DataInserted 返回 true，跳过并写入 sys_init；
	// 如果 DataInserted 返回 false，调用 InitializeData，成功后写入 sys_init。
	// 所以这里我们依赖 InitializerName 的唯一性。
	return false
}

func (i *initPatch) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 1. 补全 SysApi
	apis := []sysModel.SysApi{
		// 实例高级功能
		{ApiGroup: "instance", Method: "GET", Path: "/inst/stats", Description: "获取容器监控数据"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/webssh", Description: "WebSSH连接"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/logs", Description: "容器日志"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/file/list", Description: "获取容器文件列表"},
		{ApiGroup: "instance", Method: "GET", Path: "/inst/file/download", Description: "下载容器文件"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/file/upload", Description: "上传容器文件"},
		{ApiGroup: "instance", Method: "POST", Path: "/inst/file/delete", Description: "删除容器文件"},

		// 镜像管理
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/image/list", Description: "获取节点镜像列表"},
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/image/pull", Description: "拉取镜像"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/image/delete", Description: "删除镜像"},

		// 网络管理
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/network/list", Description: "获取网络列表"},
		{ApiGroup: "cloud", Method: "POST", Path: "/cloud/network/create", Description: "创建网络"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/network/delete", Description: "删除网络"},

		// 卷管理
		{ApiGroup: "cloud", Method: "GET", Path: "/cloud/volume/list", Description: "获取卷列表"},
		{ApiGroup: "cloud", Method: "POST", Path: "/cloud/volume/create", Description: "创建卷"},
		{ApiGroup: "cloud", Method: "DELETE", Path: "/cloud/volume/delete", Description: "删除卷"},

		// 微调任务
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

	for _, api := range apis {
		if err := db.FirstOrCreate(&api, sysModel.SysApi{Path: api.Path, Method: api.Method}).Error; err != nil {
			return ctx, errors.Wrapf(err, "补全API失败: %v %v", api.Path, api.Method)
		}
	}

	// 2. 补全 Casbin 规则 (为超级管理员 888 添加权限)
	casbinRules := []adapter.CasbinRule{}
	for _, api := range apis {
		casbinRules = append(casbinRules, adapter.CasbinRule{
			Ptype: "p",
			V0:    "888",
			V1:    api.Path,
			V2:    api.Method,
		})
	}

	// 批量插入 Casbin 规则
	// GORM Adapter 的 CasbinRule 没有唯一键约束（除了 ID），所以 FirstOrCreate 可能比较慢
	// 但为了不重复，我们还是检查一下
	for _, rule := range casbinRules {
		var count int64
		db.Model(&adapter.CasbinRule{}).Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", rule.Ptype, rule.V0, rule.V1, rule.V2).Count(&count)
		if count == 0 {
			if err := db.Create(&rule).Error; err != nil {
				return ctx, errors.Wrapf(err, "补全Casbin规则失败: %v", rule)
			}
		}
	}

	return ctx, nil
}
