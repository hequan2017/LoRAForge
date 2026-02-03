package cloud

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FineTuneTask 模型微调任务
type FineTuneTask struct {
	global.GVA_MODEL

	// ========== 基础信息 ==========
	TaskName        *string `json:"taskName" form:"taskName" gorm:"comment:任务名称;column:task_name;type:varchar(255);index"`     // 任务名称
	TaskDescription *string `json:"taskDescription" form:"taskDescription" gorm:"comment:任务描述;column:task_description;type:text"` // 任务描述
	UserID          *int64  `json:"userId" form:"userId" gorm:"comment:用户ID;column:user_id;index"`                               // 用户ID
	NodeID          *int64  `json:"nodeId" form:"nodeId" gorm:"comment:计算节点ID;column:node_id;index"`                           // 计算节点ID
	InstanceID      *int64  `json:"instanceId" form:"instanceId" gorm:"comment:关联实例ID;column:instance_id"`                    // 关联的实例ID
	ContainerID     *string `json:"containerId" form:"containerId" gorm:"comment:容器ID;column:container_id;type:varchar(64)"`   // Docker容器ID
	TaskStatus      *string `json:"taskStatus" form:"taskStatus" gorm:"comment:任务状态;column:task_status;type:varchar(32);index"` // 任务状态: pending, running, completed, failed, stopped, cancelled
	Progress        *float64 `json:"progress" form:"progress" gorm:"comment:训练进度(0-100);column:progress"`                     // 训练进度百分比
	CurrentStep     *int64  `json:"currentStep" form:"currentStep" gorm:"comment:当前训练步数;column:current_step"`                // 当前训练步数
	TotalSteps      *int64  `json:"totalSteps" form:"totalSteps" gorm:"comment:总训练步数;column:total_steps"`                    // 总训练步数

	// ========== 模型配置 ==========
	ModelType       *string `json:"modelType" form:"modelType" gorm:"comment:模型类型;column:model_type;type:varchar(32)"`     // 模型类型: llm, wa, audio
	ModelPath       *string `json:"modelPath" form:"modelPath" gorm:"comment:模型路径/ID;column:model_path;type:text"`         // 模型路径或ModelScope ID
	LoraTargetModules *string `json:"loraTargetModules" form:"loraTargetModules" gorm:"comment:LoRA目标模块;column:lora_target_modules;type:text"` // LoRA目标模块，如: c_attn, q_attn
	LoraRank        *int    `json:"loraRank" form:"loraRank" gorm:"comment:LoRA秩;column:lora_rank;default:8"`                // LoRA秩
	LoraAlpha       *int    `json:"loraAlpha" form:"loraAlpha" gorm:"comment:LoRA Alpha;column:lora_alpha;default:32"`       // LoRA Alpha
	LoraDropRate    *float64 `json:"loraDropRate" form:"loraDropRate" gorm:"comment:LoRA丢弃率;column:lora_drop_rate;default:0.05"` // LoRA Dropout
	LoraModules     *string `json:"loraModules" form:"loraModules" gorm:"comment:LoRA模块;column:lora_modules;type:text"`     // LoRA模块保存路径

	// ========== 数据集配置 ==========
	TrainDataset    *string `json:"trainDataset" form:"trainDataset" gorm:"comment:训练数据集;column:train_dataset;type:text"` // 训练数据集路径/ID
	ValDataset      *string `json:"valDataset" form:"valDataset" gorm:"comment:验证数据集;column:val_dataset;type:text"`     // 验证数据集路径/ID
	TestDataset     *string `json:"testDataset" form:"testDataset" gorm:"comment:测试数据集;column:test_dataset;type:text"`  // 测试数据集路径/ID
	DatasetType     *string `json:"datasetType" form:"datasetType" gorm:"comment:数据集类型;column:dataset_type;type:varchar(32)"` // 数据集类型: custom, alpaca, sharegpt, etc.
	DatasetProb     *float64 `json:"datasetProb" form:"datasetProb" gorm:"comment:数据集采样比例;column:dataset_prob;default:1.0"` // 数据集采样比例
	NumSamples      *int64  `json:"numSamples" form:"numSamples" gorm:"comment:训练样本数;column:num_samples"`                // 训练样本数量
	MaxSamples      *int64  `json:"maxSamples" form:"maxSamples" gorm:"comment:最大样本数;column:max_samples"`                // 最大训练样本数

	// ========== 训练参数 ==========
	LearningRate    *float64 `json:"learningRate" form:"learningRate" gorm:"comment:学习率;column:learning_rate;default:0.0002"` // 学习率
	BatchSize       *int    `json:"batchSize" form:"batchSize" gorm:"comment:批处理大小;column:batch_size;default:8"`         // 批处理大小
	GradientAccSteps *int    `json:"gradientAccSteps" form:"gradientAccSteps" gorm:"comment:梯度累积步数;column:gradient_acc_steps;default:1"` // 梯度累积步数
	NumEpochs       *int    `json:"numEpochs" form:"numEpochs" gorm:"comment:训练轮数;column:num_epochs;default:1"`          // 训练轮数
	MaxSteps        *int64  `json:"maxSteps" form:"maxSteps" gorm:"comment:最大步数;column:max_steps"`                       // 最大训练步数
	WarmupRatio     *float64 `json:"warmupRatio" form:"warmupRatio" gorm:"comment:预热比例;column:warmup_ratio;default:0.1"`  // 预热比例
	Optimizer       *string `json:"optimizer" form:"optimizer" gorm:"comment:优化器;column:optimizer;type:varchar(32);default:adamw"` // 优化器: adamw, adam, lion
	WeightDecay     *float64 `json:"weightDecay" form:"weightDecay" gorm:"comment:权重衰减;column:weight_decay;default:0.01"`  // 权重衰减
	MaxGradNorm     *float64 `json:"maxGradNorm" form:"maxGradNorm" gorm:"comment:梯度裁剪;column:max_grad_norm;default:1.0"` // 梯度裁剪

	// ========== 训练环境 ==========
	GPUCount        *int    `json:"gpuCount" form:"gpuCount" gorm:"comment:GPU数量;column:gpu_count;default:1"`              // GPU数量
	GPUType         *string `json:"gpuType" form:"gpuType" gorm:"comment:GPU类型;column:gpu_type;type:varchar(32)"`        // GPU类型: A100, A800, H800, etc.
	DockerImage     *string `json:"dockerImage" form:"dockerImage" gorm:"comment:训练镜像;column:docker_image;type:varchar(255)"` // 训练使用的Docker镜像
	SwarmPath       *string `json:"swarmPath" form:"swarmPath" gorm:"comment:SWIRT路径;column:swarm_path;type:text"`        // SWIRT工具路径

	// ========== 输出配置 ==========
	OutputDir       *string `json:"outputDir" form:"outputDir" gorm:"comment:输出目录;column:output_dir;type:varchar(255)"`  // 输出目录
	CheckpointStep  *int64  `json:"checkpointStep" form:"checkpointStep" gorm:"comment:检查点保存步数;column:checkpoint_step"` // 每N步保存检查点
	OutputCheckpoints *string `json:"outputCheckpoints" form:"outputCheckpoints" gorm:"comment:检查点输出配置;column:output_checkpoints;type:text"` // 检查点输出配置
	OnlySaveModel   *bool   `json:"onlySaveModel" form:"onlySaveModel" gorm:"comment:仅保存模型;column:only_save_model;default:false"` // 仅保存最终模型
	ResumeCheckpoint *string `json:"resumeCheckpoint" form:"resumeCheckpoint" gorm:"comment:恢复检查点路径;column:resume_checkpoint;type:text"` // 从检查点恢复训练

	// ========== 高级配置 ==========
	TemplateType    *string `json:"templateType" form:"templateType" gorm:"comment:提示模板类型;column:template_type;type:varchar(64)"` // 提示模板类型
	SequenceLength  *int    `json:"sequenceLength" form:"sequenceLength" gorm:"comment:序列长度;column:sequence_length;default:2048"` // 序列长度
	Precision       *string `json:"precision" form:"precision" gorm:"comment:训练精度;column:precision;type:varchar(16);default:bf16"` // 训练精度: bf16, fp16, fp32
	QuantizationBit *int    `json:"quantizationBit" form:"quantizationBit" gorm:"comment:量化位数;column:quantization_bit"`         // 量化位数: 4, 8
	FlashAttention  *bool   `json:"flashAttention" form:"flashAttention" gorm:"comment:启用Flash Attention;column:flash_attention;default:true"` // 启用Flash Attention
	Deepspeed       *bool   `json:"deepspeed" form:"deepspeed" gorm:"comment:启用DeepSpeed;column:deepspeed;default:false"`   // 启用DeepSpeed
	DeepspeedConfig *string `json:"deepspeedConfig" form:"deepspeedConfig" gorm:"comment:DeepSpeed配置;column:deepspeed_config;type:text"` // DeepSpeed配置文件路径

	// ========== 结果信息 ==========
	ErrorMessage    *string `json:"errorMessage" form:"errorMessage" gorm:"comment:错误信息;column:error_message;type:text"`    // 错误信息
	StartTime       *int64  `json:"startTime" form:"startTime" gorm:"comment:开始时间;column:start_time"`                       // 开始时间(Unix时间戳)
	EndTime         *int64  `json:"endTime" form:"endTime" gorm:"comment:结束时间;column:end_time"`                           // 结束时间(Unix时间戳)
	Duration        *int64  `json:"duration" form:"duration" gorm:"comment:耗时(秒);column:duration"`                         // 训练耗时(秒)
	OutputModelPath *string `json:"outputModelPath" form:"outputModelPath" gorm:"comment:输出模型路径;column:output_model_path;type:text"` // 输出模型路径
	LogPath         *string `json:"logPath" form:"logPath" gorm:"comment:日志路径;column:log_path;type:text"`                 // 日志文件路径
	TensorboardPath *string `json:"tensorboardPath" form:"tensorboardPath" gorm:"comment:TensorBoard路径;column:tensorboard_path;type:text"` // TensorBoard路径
	ValidationLoss  *float64 `json:"validationLoss" form:"validationLoss" gorm:"comment:验证损失;column:validation_loss"`      // 最终验证损失

	// ========== 其他配置 ==========
	Remark          *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;type:text"`                          // 备注
	CustomConfig    *string `json:"customConfig" form:"customConfig" gorm:"type:json;comment:自定义配置(JSON);column:custom_config"` // 自定义配置(JSON格式)
}

// TableName FineTuneTask表名
func (FineTuneTask) TableName() string {
	return "finetune_tasks"
}

// FineTuneTaskSnapshot 训练快照记录（用于存储训练过程中的指标）
type FineTuneTaskSnapshot struct {
	global.GVA_MODEL

	TaskID      *uint   `json:"taskId" form:"taskId" gorm:"comment:任务ID;column:task_id;index"`                // 关联的任务ID
	Step        *int64  `json:"step" form:"step" gorm:"comment:训练步数;column:step;index"`                      // 训练步数
	Loss        *float64 `json:"loss" form:"loss" gorm:"comment:训练损失;column:loss"`                          // 训练损失
	LearningRate *float64 `json:"learningRate" form:"learningRate" gorm:"comment:当前学习率;column:learning_rate"` // 当前学习率
	TrainSpeed  *float64 `json:"trainSpeed" form:"trainSpeed" gorm:"comment:训练速度(samples/s);column:train_speed"` // 训练速度
	GPUMemory   *float64 `json:"gpuMemory" form:"gpuMemory" gorm:"comment:GPU显存占用(GB);column:gpu_memory"`  // GPU显存占用
	Timestamp   *int64  `json:"timestamp" form:"timestamp" gorm:"comment:时间戳;column:timestamp"`              // 时间戳
}

// TableName FineTuneTaskSnapshot表名
func (FineTuneTaskSnapshot) TableName() string {
	return "finetune_task_snapshots"
}
