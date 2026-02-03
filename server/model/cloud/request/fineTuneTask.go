package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// FineTuneTaskSearch 微调任务搜索结构体
type FineTuneTaskSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	TaskName       string      `json:"taskName" form:"taskName"`             // 任务名称(模糊查询)
	TaskStatus     string      `json:"taskStatus" form:"taskStatus"`         // 任务状态(精确匹配)
	ModelType      string      `json:"modelType" form:"modelType"`           // 模型类型
	UserID         *int64      `json:"userId" form:"userId"`                 // 用户ID
	NodeID         *int64      `json:"nodeId" form:"nodeId"`                 // 节点ID
	request.PageInfo
}

// FineTuneTaskCreate 创建微调任务请求
type FineTuneTaskCreate struct {
	TaskName        *string  `json:"taskName" binding:"required"`              // 任务名称(必填)
	TaskDescription *string  `json:"taskDescription"`                          // 任务描述
	NodeID          *int64   `json:"nodeId" binding:"required"`               // 计算节点ID(必填)

	// 模型配置
	ModelType        *string  `json:"modelType" binding:"required"`            // 模型类型: llm, wa, audio
	ModelPath        *string  `json:"modelPath" binding:"required"`           // 模型路径/ID
	LoraTargetModules *string `json:"loraTargetModules"`                      // LoRA目标模块
	LoraRank         *int    `json:"loraRank"`                                // LoRA秩
	LoraAlpha        *int    `json:"loraAlpha"`                               // LoRA Alpha
	LoraDropRate     *float64 `json:"loraDropRate"`                           // LoRA丢弃率

	// 数据集配置
	TrainDataset *string  `json:"trainDataset" binding:"required"`           // 训练数据集(必填)
	ValDataset   *string  `json:"valDataset"`                                // 验证数据集
	DatasetType  *string  `json:"datasetType"`                               // 数据集类型
	DatasetProb  *float64 `json:"datasetProb"`                               // 数据集采样比例
	MaxSamples   *int64   `json:"maxSamples"`                                // 最大样本数

	// 训练参数
	LearningRate     *float64 `json:"learningRate"`                           // 学习率
	BatchSize        *int    `json:"batchSize"`                              // 批处理大小
	GradientAccSteps *int    `json:"gradientAccSteps"`                       // 梯度累积步数
	NumEpochs        *int    `json:"numEpochs"`                              // 训练轮数
	MaxSteps         *int64  `json:"maxSteps"`                               // 最大训练步数
	WarmupRatio      *float64 `json:"warmupRatio"`                           // 预热比例
	Optimizer        *string `json:"optimizer"`                              // 优化器
	WeightDecay      *float64 `json:"weightDecay"`                           // 权重衰减
	MaxGradNorm      *float64 `json:"maxGradNorm"`                           // 梯度裁剪

	// 训练环境
	GPUCount    *int    `json:"gpuCount"`                                   // GPU数量
	GPUType     *string `json:"gpuType"`                                    // GPU类型
	DockerImage *string `json:"dockerImage"`                                // Docker镜像

	// 输出配置
	OutputDir        *string `json:"outputDir"`                             // 输出目录
	CheckpointStep   *int64  `json:"checkpointStep"`                        // 检查点保存步数
	OnlySaveModel    *bool   `json:"onlySaveModel"`                         // 仅保存最终模型
	ResumeCheckpoint *string `json:"resumeCheckpoint"`                      // 恢复检查点路径

	// 高级配置
	TemplateType    *string `json:"templateType"`                           // 提示模板类型
	SequenceLength  *int    `json:"sequenceLength"`                         // 序列长度
	Precision       *string `json:"precision"`                              // 训练精度
	QuantizationBit *int    `json:"quantizationBit"`                        // 量化位数
	FlashAttention  *bool   `json:"flashAttention"`                         // 启用Flash Attention
	Deepspeed        *bool   `json:"deepspeed"`                             // 启用DeepSpeed
	DeepspeedConfig *string `json:"deepspeedConfig"`                        // DeepSpeed配置

	CustomConfig *string `json:"customConfig"`                             // 自定义配置(JSON)
	Remark       *string `json:"remark"`                                   // 备注
}

// FineTuneTaskUpdate 更新微调任务请求
type FineTuneTaskUpdate struct {
	ID              *uint   `json:"id" binding:"required"`
	TaskName        *string `json:"taskName"`
	TaskDescription *string `json:"taskDescription"`
	TaskStatus      *string `json:"taskStatus"`                             // 仅允许更新到特定状态(如停止、取消)
	Progress        *float64 `json:"progress"`
	CurrentStep     *int64  `json:"currentStep"`
	TotalSteps      *int64  `json:"totalSteps"`
	ErrorMessage    *string `json:"errorMessage"`
	OutputModelPath *string `json:"outputModelPath"`
	ValidationLoss  *float64 `json:"validationLoss"`
	Remark          *string `json:"remark"`
}

// FineTuneTaskStart 启动微调任务请求
type FineTuneTaskStart struct {
	ID *uint `json:"id" binding:"required"`                                 // 任务ID
}

// FineTuneTaskStop 停止微调任务请求
type FineTuneTaskStop struct {
	ID     *uint   `json:"id" binding:"required"`                           // 任务ID
	Reason *string `json:"reason"`                                          // 停止原因
}

// FineTuneTaskDelete 删除微调任务请求
type FineTuneTaskDelete struct {
	ID *uint `json:"id" binding:"required"`                                // 任务ID
}

// FineTuneTaskBatchDelete 批量删除微调任务请求
type FineTuneTaskBatchDelete struct {
	IDs []uint `json:"ids" binding:"required"`                             // 任务ID列表
}

// FineTuneTaskLogRequest 获取任务日志请求
type FineTuneTaskLogRequest struct {
	ID     *uint  `json:"id" binding:"required"`                           // 任务ID
	Offset *int   `json:"offset"`                                          // 日志偏移量
	Limit  *int   `json:"limit"`                                           // 日志行数限制
	Follow *bool  `json:"follow"`                                          // 是否持续监听
}

// FineTuneTaskSnapshotCreate 创建训练快照请求
type FineTuneTaskSnapshotCreate struct {
	TaskID       *uint    `json:"taskId" binding:"required"`               // 任务ID
	Step         *int64   `json:"step" binding:"required"`                 // 训练步数
	Loss         *float64 `json:"loss"`                                    // 训练损失
	LearningRate *float64 `json:"learningRate"`                            // 当前学习率
	TrainSpeed   *float64 `json:"trainSpeed"`                              // 训练速度
	GPUMemory    *float64 `json:"gpuMemory"`                               // GPU显存占用
}

// FineTuneTaskSnapshotSearch 训练快照搜索请求
type FineTuneTaskSnapshotSearch struct {
	TaskID *uint `json:"taskId" binding:"required"`                        // 任务ID
	StartTime *int64 `json:"startTime"`                                     // 开始时间戳
	EndTime   *int64 `json:"endTime"`                                       // 结束时间戳
	request.PageInfo
}

// FineTuneTaskExecuteRequest 执行训练请求（调用SWIRT API）
type FineTuneTaskExecuteRequest struct {
	TaskID      *uint                   `json:"taskId" binding:"required"` // 任务ID
	Environment map[string]string       `json:"environment"`                // 环境变量
	Mounts      []string                `json:"mounts"`                     // 挂载配置
}

// FineTuneTaskMonitorRequest 监控任务状态请求
type FineTuneTaskMonitorRequest struct {
	TaskID *uint `json:"taskId" binding:"required"`                       // 任务ID
}

// FineTuneTaskCancelRequest 取消任务请求
type FineTuneTaskCancelRequest struct {
	ID     *uint   `json:"id" binding:"required"`                         // 任务ID
	Reason *string `json:"reason"`                                        // 取消原因
}

// FineTuneTaskResumeRequest 恢复任务请求
type FineTuneTaskResumeRequest struct {
	ID               *uint   `json:"id" binding:"required"`               // 任务ID
	ResumeCheckpoint *string `json:"resumeCheckpoint"`                     // 检查点路径
}

// FineTuneTaskExportRequest 导出模型请求
type FineTuneTaskExportRequest struct {
	ID           *uint   `json:"id" binding:"required"`                   // 任务ID
	ExportFormat *string `json:"exportFormat"`                            // 导出格式: merge, vllm, awq, gptq
	ExportPath   *string `json:"exportPath"`                              // 导出路径
	QuantizeBit  *int    `json:"quantizeBit"`                             // 量化位数(仅用于AWQ/GPTQ)
}
