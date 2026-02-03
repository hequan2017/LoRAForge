package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
)

type FineTuneTaskService struct{}

// CreateFineTuneTask 创建微调任务
func (s *FineTuneTaskService) CreateFineTuneTask(ctx context.Context, task *cloud.FineTuneTask) (err error) {
	global.GVA_LOG.Info("开始创建微调任务", zap.Any("task", task))

	// 1. 验证必填字段
	if task.TaskName == nil || *task.TaskName == "" {
		return fmt.Errorf("任务名称不能为空")
	}
	if task.NodeID == nil || *task.NodeID == 0 {
		return fmt.Errorf("节点ID不能为空")
	}
	if task.ModelPath == nil || *task.ModelPath == "" {
		return fmt.Errorf("模型路径不能为空")
	}
	if task.TrainDataset == nil || *task.TrainDataset == "" {
		return fmt.Errorf("训练数据集不能为空")
	}

	// 2. 检查节点是否可用
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 3. 检查当前运行中的任务数量
	var runningCount int64
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Where("task_status = ? AND deleted_at IS NULL", "running").
		Count(&runningCount)
	if runningCount >= 8 {
		return fmt.Errorf("系统运行中的任务已达到上限(8个)，请稍后再试")
	}

	// 4. 设置默认值
	if task.TaskStatus == nil {
		status := "pending"
		task.TaskStatus = &status
	}
	if task.Progress == nil {
		progress := 0.0
		task.Progress = &progress
	}
	if task.CurrentStep == nil {
		step := int64(0)
		task.CurrentStep = &step
	}
	if task.LoraRank == nil {
		rank := 8
		task.LoraRank = &rank
	}
	if task.LoraAlpha == nil {
		alpha := 32
		task.LoraAlpha = &alpha
	}
	if task.LoraDropRate == nil {
		dropRate := 0.05
		task.LoraDropRate = &dropRate
	}
	if task.LearningRate == nil {
		lr := 0.0002
		task.LearningRate = &lr
	}
	if task.BatchSize == nil {
		batchSize := 8
		task.BatchSize = &batchSize
	}
	if task.NumEpochs == nil {
		epochs := 1
		task.NumEpochs = &epochs
	}
	if task.WarmupRatio == nil {
		warmup := 0.1
		task.WarmupRatio = &warmup
	}
	if task.Optimizer == nil {
		optimizer := "adamw"
		task.Optimizer = &optimizer
	}
	if task.WeightDecay == nil {
		decay := 0.01
		task.WeightDecay = &decay
	}
	if task.MaxGradNorm == nil {
		norm := 1.0
		task.MaxGradNorm = &norm
	}
	if task.GPUCount == nil {
		gpuCount := 1
		task.GPUCount = &gpuCount
	}
	if task.Precision == nil {
		precision := "bf16"
		task.Precision = &precision
	}
	if task.FlashAttention == nil {
		flashAttention := true
		task.FlashAttention = &flashAttention
	}
	if task.OnlySaveModel == nil {
		onlySaveModel := false
		task.OnlySaveModel = &onlySaveModel
	}
	if task.DatasetProb == nil {
		prob := 1.0
		task.DatasetProb = &prob
	}

	// 5. 设置输出目录
	if task.OutputDir == nil || *task.OutputDir == "" {
		outputDir := fmt.Sprintf("/data/output/finetune/%d", time.Now().Unix())
		task.OutputDir = &outputDir
	}

	// 6. 创建数据库记录
	if err := global.GVA_DB.Create(task).Error; err != nil {
		global.GVA_LOG.Error("创建微调任务记录失败", zap.Error(err))
		return fmt.Errorf("创建微调任务记录失败: %v", err)
	}

	global.GVA_LOG.Info("微调任务创建成功", zap.Uint("taskId", task.ID))
	return nil
}

// StartFineTuneTask 启动微调任务
func (s *FineTuneTaskService) StartFineTuneTask(ctx context.Context, taskID uint) (err error) {
	global.GVA_LOG.Info("开始启动微调任务", zap.Uint("taskId", taskID))

	// 1. 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("未找到微调任务: %v", err)
	}

	// 2. 检查任务状态
	if *task.TaskStatus != "pending" {
		return fmt.Errorf("任务状态不允许启动，当前状态: %s", *task.TaskStatus)
	}

	// 3. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 4. 检查用户是否有正在运行的任务
	var userRunningCount int64
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Where("user_id = ? AND task_status = ? AND id != ? AND deleted_at IS NULL", task.UserID, "running", taskID).
		Count(&userRunningCount)
	if userRunningCount > 0 {
		return fmt.Errorf("您已有正在运行的任务，请等待完成后再启动新任务")
	}

	// 5. 创建 Docker 客户端
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 6. 准备训练镜像
	imageName := "modelscope-registry.cn-hangzhou.cr.aliyuncs.com/modelscope-repo/modelscope:ubuntu22.04-cuda12.9.1-py311-torch2.8.0-vllm0.11.0-modelscope1.32.0-swift3.11.3" // 默认ModelScope训练镜像
	if task.DockerImage != nil && *task.DockerImage != "" {
		imageName = *task.DockerImage
	}

	global.GVA_LOG.Info("准备拉取训练镜像", zap.String("image", imageName))
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err == nil {
		defer reader.Close()
		io.Copy(io.Discard, reader)
		global.GVA_LOG.Info("训练镜像拉取成功")
	} else {
		global.GVA_LOG.Warn("训练镜像拉取失败或已存在", zap.Error(err))
	}

	// 7. 构建SWIRT训练命令
	trainCommand := s.buildSwiftCommand(&task)

	// 8. 准备容器配置
	containerName := fmt.Sprintf("finetune-%d-%d", task.ID, time.Now().Unix())
	config := &container.Config{
		Image: imageName,
		Cmd:   strings.Fields(trainCommand),
		Tty:   false,
		Env:   s.buildEnvVars(&task),
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "on-failure", MaximumRetryCount: 3},
		ShmSize:       8 * 1024 * 1024 * 1024, // 8GB shared memory
	}

	// GPU 配置
	if task.GPUCount != nil && *task.GPUCount > 0 {
		hostConfig.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Count:        *task.GPUCount,
				Capabilities: [][]string{{"gpu"}},
			},
		}
	}

	// 挂载配置
	var binds []string
	// 挂载输出目录
	outputBind := fmt.Sprintf("finetune-output-%d:%s", task.ID, *task.OutputDir)
	binds = append(binds, outputBind)

	// 挂载数据集目录 - 仅挂载绝对路径（ModelScope数据集ID会由Swift自动下载）
	if task.TrainDataset != nil && *task.TrainDataset != "" {
		if strings.HasPrefix(*task.TrainDataset, "/") {
			// 本地绝对路径，需要挂载
			datasetBind := fmt.Sprintf("%s:%s", *task.TrainDataset, *task.TrainDataset)
			binds = append(binds, datasetBind)
		}
		// ModelScope ID（如: AI-ModelScope/xxx）不需要挂载，Swift会自动下载
	}

	// 挂载模型目录 - 仅挂载绝对路径
	if task.ModelPath != nil && *task.ModelPath != "" {
		if strings.HasPrefix(*task.ModelPath, "/") {
			// 本地绝对路径，需要挂载
			modelBind := fmt.Sprintf("%s:%s", *task.ModelPath, *task.ModelPath)
			binds = append(binds, modelBind)
		}
		// ModelScope ID（如: qwen/Qwen-7B）不需要挂载，Swift会自动下载
	}

	// 挂载模型缓存目录（ModelScope下载的模型会缓存到这里）
	binds = append(binds, "modelscope-cache:/data/models")
	// 挂载HuggingFace缓存目录
	binds = append(binds, "hf-cache:/data/huggingface")

	hostConfig.Binds = binds

	// 9. 创建容器
	global.GVA_LOG.Info("准备创建训练容器",
		zap.String("name", containerName),
		zap.String("command", trainCommand),
	)
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		global.GVA_LOG.Error("创建训练容器失败", zap.Error(err))
		return fmt.Errorf("创建训练容器失败: %v", err)
	}

	// 10. 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		global.GVA_LOG.Error("启动训练容器失败", zap.Error(err))
		return fmt.Errorf("启动训练容器失败: %v", err)
	}

	// 11. 更新任务状态
	shortID := resp.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}
	task.ContainerID = &shortID
	status := "running"
	task.TaskStatus = &status
	task.StartTime = new(int64)
	*task.StartTime = time.Now().Unix()

	if err := global.GVA_DB.Save(&task).Error; err != nil {
		global.GVA_LOG.Error("更新任务状态失败", zap.Error(err))
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	global.GVA_LOG.Info("微调任务启动成功", zap.Uint("taskId", taskID), zap.String("containerId", shortID))
	return nil
}

// buildSwiftCommand 构建Swift训练命令
func (s *FineTuneTaskService) buildSwiftCommand(task *cloud.FineTuneTask) string {
	var cmdParts []string

	// 基础命令 - 使用 sft (Supervised Fine-Tuning) 进行监督微调
	cmdParts = append(cmdParts, "swift", "sft")

	// 模型配置
	cmdParts = append(cmdParts, fmt.Sprintf("--model_type=%s", *task.ModelType))
	cmdParts = append(cmdParts, fmt.Sprintf("--model_id_or_path=%s", *task.ModelPath))

	// 微调类型 - 使用 LoRA
	cmdParts = append(cmdParts, "--sft_type=lora")

	// 数据集配置
	cmdParts = append(cmdParts, fmt.Sprintf("--dataset=%s", *task.TrainDataset))
	if task.ValDataset != nil && *task.ValDataset != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("--val_dataset=%s", *task.ValDataset))
	}

	// LoRA配置
	if task.LoraTargetModules != nil && *task.LoraTargetModules != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("--lora_target_modules=%s", *task.LoraTargetModules))
	}
	cmdParts = append(cmdParts, fmt.Sprintf("--lora_rank=%d", *task.LoraRank))
	cmdParts = append(cmdParts, fmt.Sprintf("--lora_alpha=%d", *task.LoraAlpha))
	if task.LoraDropRate != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--lora_dropout=%f", *task.LoraDropRate))
	}

	// 训练参数
	cmdParts = append(cmdParts, fmt.Sprintf("--learning_rate=%f", *task.LearningRate))
	cmdParts = append(cmdParts, fmt.Sprintf("--batch_size=%d", *task.BatchSize))
	if task.GradientAccSteps != nil && *task.GradientAccSteps > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--gradient_accumulation_steps=%d", *task.GradientAccSteps))
	}
	if task.NumEpochs != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--num_train_epochs=%d", *task.NumEpochs))
	}
	if task.MaxSteps != nil && *task.MaxSteps > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--max_steps=%d", *task.MaxSteps))
	}
	if task.WarmupRatio != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--warmup_ratio=%f", *task.WarmupRatio))
	}
	if task.Optimizer != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--optim=%s", *task.Optimizer))
	}
	if task.WeightDecay != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--weight_decay=%f", *task.WeightDecay))
	}
	if task.MaxGradNorm != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--max_grad_norm=%f", *task.MaxGradNorm))
	}

	// 输出配置
	cmdParts = append(cmdParts, fmt.Sprintf("--output_dir=%s", *task.OutputDir))
	if task.CheckpointStep != nil && *task.CheckpointStep > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--save_steps=%d", *task.CheckpointStep))
	}
	if task.OnlySaveModel != nil && *task.OnlySaveModel {
		cmdParts = append(cmdParts, "--only_save_model=true")
	}

	// 高级配置
	if task.TemplateType != nil && *task.TemplateType != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("--template_type=%s", *task.TemplateType))
	}
	if task.SequenceLength != nil {
		cmdParts = append(cmdParts, fmt.Sprintf("--max_length=%d", *task.SequenceLength))
	}
	if task.Precision != nil {
		precisionMap := map[string]string{
			"bf16": "--bf16=true",
			"fp16": "--fp16=true",
			"fp32": "--fp16=false",
		}
		if cmd, ok := precisionMap[*task.Precision]; ok {
			cmdParts = append(cmdParts, cmd)
		}
	}
	if task.QuantizationBit != nil && *task.QuantizationBit > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--quantization_bit=%d", *task.QuantizationBit))
	}
	if task.FlashAttention != nil && *task.FlashAttention {
		cmdParts = append(cmdParts, "--flash_attention=true")
	}
	if task.Deepspeed != nil && *task.Deepspeed {
		cmdParts = append(cmdParts, "--deepspeed=true")
		if task.DeepspeedConfig != nil && *task.DeepspeedConfig != "" {
			cmdParts = append(cmdParts, fmt.Sprintf("--deepspeed_config=%s", *task.DeepspeedConfig))
		}
	}

	// 其他参数
	if task.DatasetProb != nil && *task.DatasetProb < 1.0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--dataset_prob=%f", *task.DatasetProb))
	}
	if task.MaxSamples != nil && *task.MaxSamples > 0 {
		cmdParts = append(cmdParts, fmt.Sprintf("--max_samples=%d", *task.MaxSamples))
	}
	if task.ResumeCheckpoint != nil && *task.ResumeCheckpoint != "" {
		cmdParts = append(cmdParts, fmt.Sprintf("--resume_from_checkpoint=%s", *task.ResumeCheckpoint))
	}

	return strings.Join(cmdParts, " ")
}

// buildEnvVars 构建环境变量
func (s *FineTuneTaskService) buildEnvVars(task *cloud.FineTuneTask) []string {
	envs := []string{
		"SWIFT_OUTPUT_DIR=" + *task.OutputDir,
		"SWIFT_LOG_LEVEL=INFO",
	}

	// 添加 ModelScope 相关配置
	envs = append(envs, "MODELSCOPE_CACHE=/data/models")
	envs = append(envs, "HF_HOME=/data/huggingface")

	return envs
}

// StopFineTuneTask 停止微调任务
func (s *FineTuneTaskService) StopFineTuneTask(ctx context.Context, taskID uint, reason string) (err error) {
	global.GVA_LOG.Info("开始停止微调任务", zap.Uint("taskId", taskID))

	// 1. 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("未找到微调任务: %v", err)
	}

	// 2. 检查任务状态
	if *task.TaskStatus != "running" {
		return fmt.Errorf("任务未在运行中，当前状态: %s", *task.TaskStatus)
	}

	// 3. 停止容器
	if task.ContainerID != nil && *task.ContainerID != "" {
		var node cloud.ComputeNode
		if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
			return fmt.Errorf("未找到计算节点: %v", err)
		}

		cli, err := CreateDockerClient(node)
		if err != nil {
			return fmt.Errorf("创建 Docker 客户端失败: %v", err)
		}
		defer cli.Close()

		timeout := 30
		if err := cli.ContainerStop(ctx, *task.ContainerID, container.StopOptions{Timeout: &timeout}); err != nil {
			global.GVA_LOG.Error("停止容器失败", zap.Error(err))
			return fmt.Errorf("停止容器失败: %v", err)
		}
	}

	// 4. 更新任务状态
	status := "stopped"
	task.TaskStatus = &status
	if reason != "" {
		task.ErrorMessage = &reason
	}
	now := time.Now().Unix()
	task.EndTime = &now
	if task.StartTime != nil {
		duration := now - *task.StartTime
		task.Duration = &duration
	}

	if err := global.GVA_DB.Save(&task).Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	global.GVA_LOG.Info("微调任务已停止", zap.Uint("taskId", taskID))
	return nil
}

// RestartFineTuneTask 重启失败的微调任务
func (s *FineTuneTaskService) RestartFineTuneTask(ctx context.Context, taskID uint) (err error) {
	global.GVA_LOG.Info("开始重启微调任务", zap.Uint("taskId", taskID))

	// 1. 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("未找到微调任务: %v", err)
	}

	// 2. 检查任务状态 - 只允许重启失败、已停止或已取消的任务
	validStatuses := map[string]bool{
		"failed":   true,
		"stopped":  true,
		"cancelled": true,
	}
	if !validStatuses[*task.TaskStatus] {
		return fmt.Errorf("任务状态不允许重启，当前状态: %s，只允许重启失败/已停止/已取消的任务", *task.TaskStatus)
	}

	// 3. 清理旧容器（如果存在）
	if task.ContainerID != nil && *task.ContainerID != "" {
		var node cloud.ComputeNode
		if err := global.GVA_DB.First(&node, task.NodeID).Error; err == nil {
			cli, err := CreateDockerClient(node)
			if err == nil {
				// 尝试停止并删除旧容器
				cli.ContainerStop(ctx, *task.ContainerID, container.StopOptions{Timeout: nil})
				cli.ContainerRemove(ctx, *task.ContainerID, container.RemoveOptions{Force: true})
				cli.Close()
				global.GVA_LOG.Info("已清理旧容器", zap.String("containerId", *task.ContainerID))
			}
		}
	}

	// 4. 重置任务状态
	status := "pending"
	task.TaskStatus = &status
	task.ContainerID = nil
	task.ErrorMessage = nil
	task.Progress = new(float64)
	*task.Progress = 0
	task.CurrentStep = new(int64)
	*task.CurrentStep = 0
	task.StartTime = nil
	task.EndTime = nil
	task.Duration = nil
	task.ValidationLoss = nil

	if err := global.GVA_DB.Save(&task).Error; err != nil {
		global.GVA_LOG.Error("更新任务状态失败", zap.Error(err))
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	// 5. 启动任务
	if err := s.StartFineTuneTask(ctx, taskID); err != nil {
		return fmt.Errorf("启动任务失败: %v", err)
	}

	global.GVA_LOG.Info("微调任务重启成功", zap.Uint("taskId", taskID))
	return nil
}

// DeleteFineTuneTask 删除微调任务
func (s *FineTuneTaskService) DeleteFineTuneTask(ctx context.Context, taskID uint) (err error) {
	// 1. 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("未找到微调任务: %v", err)
	}

	// 2. 如果任务正在运行，先停止
	if *task.TaskStatus == "running" {
		if err := s.StopFineTuneTask(ctx, taskID, "任务被删除"); err != nil {
			global.GVA_LOG.Warn("停止运行中的任务失败", zap.Error(err))
		}
	}

	// 3. 清理容器
	if task.ContainerID != nil && *task.ContainerID != "" {
		var node cloud.ComputeNode
		if err := global.GVA_DB.First(&node, task.NodeID).Error; err == nil {
			cli, err := CreateDockerClient(node)
			if err == nil {
				defer cli.Close()
				cli.ContainerRemove(ctx, *task.ContainerID, container.RemoveOptions{Force: true})
			}
		}
	}

	// 4. 删除数据库记录
	if err := global.GVA_DB.Delete(&task).Error; err != nil {
		return fmt.Errorf("删除任务记录失败: %v", err)
	}

	global.GVA_LOG.Info("微调任务已删除", zap.Uint("taskId", taskID))
	return nil
}

// DeleteFineTuneTaskByIds 批量删除微调任务
func (s *FineTuneTaskService) DeleteFineTuneTaskByIds(ctx context.Context, taskIDs []uint) (err error) {
	for _, id := range taskIDs {
		if err := s.DeleteFineTuneTask(ctx, id); err != nil {
			global.GVA_LOG.Error("删除任务失败", zap.Uint("taskId", id), zap.Error(err))
			return err
		}
	}
	return nil
}

// UpdateFineTuneTask 更新微调任务
func (s *FineTuneTaskService) UpdateFineTuneTask(ctx context.Context, task *cloud.FineTuneTask) (err error) {
	err = global.GVA_DB.Model(&cloud.FineTuneTask{}).Where("id = ?", task.ID).Updates(&task).Error
	return err
}

// GetFineTuneTask 根据ID获取微调任务
func (s *FineTuneTaskService) GetFineTuneTask(ctx context.Context, taskID uint) (task cloud.FineTuneTask, err error) {
	err = global.GVA_DB.Where("id = ?", taskID).First(&task).Error
	return
}

// GetFineTuneTaskList 分页获取微调任务列表
func (s *FineTuneTaskService) GetFineTuneTaskList(ctx context.Context, info cloudReq.FineTuneTaskSearch) (list []cloud.FineTuneTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&cloud.FineTuneTask{})

	// 搜索条件
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.TaskName != "" {
		db = db.Where("task_name LIKE ?", "%"+info.TaskName+"%")
	}
	if info.TaskStatus != "" {
		db = db.Where("task_status = ?", info.TaskStatus)
	}
	if info.ModelType != "" {
		db = db.Where("model_type = ?", info.ModelType)
	}
	if info.UserID != nil {
		db = db.Where("user_id = ?", *info.UserID)
	}
	if info.NodeID != nil {
		db = db.Where("node_id = ?", *info.NodeID)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 按创建时间倒序排列
	err = db.Order("created_at DESC").Find(&list).Error
	return
}

// GetFineTuneTaskDataSource 获取微调任务数据源
func (s *FineTuneTaskService) GetFineTuneTaskDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	// 获取计算节点列表
	nodes := make([]map[string]any, 0)
	global.GVA_DB.Table("compute_nodes").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&nodes)
	res["nodeId"] = nodes

	// 获取用户列表
	users := make([]map[string]any, 0)
	global.GVA_DB.Table("sys_users").Where("deleted_at IS NULL").Select("nick_name as label,id as value").Scan(&users)
	res["userId"] = users

	return
}

// UpdateTaskStatus 更新任务状态（供定时任务调用）
func (s *FineTuneTaskService) UpdateTaskStatus(ctx context.Context) (err error) {
	// 获取所有运行中的任务
	var runningTasks []cloud.FineTuneTask
	if err := global.GVA_DB.Where("task_status = ? AND container_id IS NOT NULL AND container_id != ''", "running").Find(&runningTasks).Error; err != nil {
		return err
	}

	for _, task := range runningTasks {
		// 获取节点信息
		var node cloud.ComputeNode
		if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
			continue
		}

		cli, err := CreateDockerClient(node)
		if err != nil {
			continue
		}

		// 检查容器状态
		containerJSON, err := cli.ContainerInspect(ctx, *task.ContainerID)
		cli.Close()

		if err != nil {
			// 容器不存在或出错，标记任务为失败
			status := "failed"
			errorMsg := "容器不存在或已被删除"
			global.GVA_DB.Model(&task).Updates(map[string]any{
				"task_status":   status,
				"error_message": errorMsg,
			})
			continue
		}

		// 检查容器是否已退出
		if !containerJSON.State.Running {
			exitCode := containerJSON.State.ExitCode
			if exitCode == 0 {
				status := "completed"
				now := time.Now().Unix()
				global.GVA_DB.Model(&task).Updates(map[string]any{
					"task_status": status,
					"progress":    100.0,
					"end_time":    now,
				})
			} else {
				status := "failed"
				errorMsg := fmt.Sprintf("容器异常退出，退出码: %d", exitCode)
				global.GVA_DB.Model(&task).Updates(map[string]any{
					"task_status":   status,
					"error_message": errorMsg,
				})
			}
		}
	}

	return nil
}

// GetTaskLogs 获取任务日志
func (s *FineTuneTaskService) GetTaskLogs(ctx context.Context, taskID uint, offset, limit int) (logs []string, err error) {
	// 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("未找到微调任务: %v", err)
	}

	if task.ContainerID == nil || *task.ContainerID == "" {
		return nil, fmt.Errorf("任务未关联容器")
	}

	// 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 获取容器日志
	reader, err := cli.ContainerLogs(ctx, *task.ContainerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", limit),
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器日志失败: %v", err)
	}
	defer reader.Close()

	// 读取日志
	logBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取日志失败: %v", err)
	}

	logLines := strings.Split(string(logBytes), "\n")
	if offset >= len(logLines) {
		return []string{}, nil
	}

	if offset+limit < len(logLines) {
		return logLines[offset : offset+limit], nil
	}
	return logLines[offset:], nil
}

// CreateSnapshot 创建训练快照
func (s *FineTuneTaskService) CreateSnapshot(ctx context.Context, snapshot *cloud.FineTuneTaskSnapshot) error {
	snapshot.Timestamp = new(int64)
	*snapshot.Timestamp = time.Now().Unix()
	return global.GVA_DB.Create(snapshot).Error
}

// GetSnapshots 获取任务训练快照列表
func (s *FineTuneTaskService) GetSnapshots(ctx context.Context, taskID uint) (snapshots []cloud.FineTuneTaskSnapshot, err error) {
	err = global.GVA_DB.Where("task_id = ?", taskID).Order("step ASC").Find(&snapshots).Error
	return
}

// GetMetrics 获取任务训练指标
func (s *FineTuneTaskService) GetMetrics(ctx context.Context, taskID uint) (metrics map[string]any, err error) {
	metrics = make(map[string]any)

	// 获取任务基本信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	metrics["taskStatus"] = task.TaskStatus
	metrics["progress"] = task.Progress
	metrics["currentStep"] = task.CurrentStep
	metrics["totalSteps"] = task.TotalSteps
	metrics["validationLoss"] = task.ValidationLoss

	// 获取快照数据
	snapshots, err := s.GetSnapshots(ctx, taskID)
	if err == nil && len(snapshots) > 0 {
		steps := make([]int64, 0, len(snapshots))
		losses := make([]float64, 0, len(snapshots))
		learningRates := make([]float64, 0, len(snapshots))
		trainSpeeds := make([]float64, 0, len(snapshots))
		gpuMemories := make([]float64, 0, len(snapshots))

		for _, snap := range snapshots {
			if snap.Step != nil {
				steps = append(steps, *snap.Step)
			}
			if snap.Loss != nil {
				losses = append(losses, *snap.Loss)
			}
			if snap.LearningRate != nil {
				learningRates = append(learningRates, *snap.LearningRate)
			}
			if snap.TrainSpeed != nil {
				trainSpeeds = append(trainSpeeds, *snap.TrainSpeed)
			}
			if snap.GPUMemory != nil {
				gpuMemories = append(gpuMemories, *snap.GPUMemory)
			}
		}

		metrics["steps"] = steps
		metrics["losses"] = losses
		metrics["learningRates"] = learningRates
		metrics["trainSpeeds"] = trainSpeeds
		metrics["gpuMemories"] = gpuMemories
	}

	return metrics, nil
}

// ExportModel 导出模型
func (s *FineTuneTaskService) ExportModel(ctx context.Context, taskID uint, exportFormat, exportPath string, quantizeBit int) error {
	// 获取任务信息
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("未找到微调任务: %v", err)
	}

	if *task.TaskStatus != "completed" {
		return fmt.Errorf("只能导出已完成的模型")
	}

	// 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, task.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 构建导出命令
	var exportCmd []string
	outputPath := exportPath
	if outputPath == "" {
		outputPath = filepath.Join(*task.OutputDir, "exported")
	}

	switch exportFormat {
	case "merge":
		exportCmd = []string{
			"swift", "export",
			fmt.Sprintf("--model_type=%s", *task.ModelType),
			fmt.Sprintf("--model_id_or_path=%s", *task.ModelPath),
			fmt.Sprintf("--ckpt_dir=%s", *task.OutputDir),
			"--merge_lora=true",
			fmt.Sprintf("--output_dir=%s", outputPath),
		}
	case "vllm":
		exportCmd = []string{
			"swift", "export",
			fmt.Sprintf("--model_type=%s", *task.ModelType),
			fmt.Sprintf("--model_id_or_path=%s", *task.ModelPath),
			fmt.Sprintf("--ckpt_dir=%s", *task.OutputDir),
			"--merge_lora=true",
			fmt.Sprintf("--output_dir=%s", outputPath),
			"--to_vllm=true",
		}
	case "awq":
		exportCmd = []string{
			"swift", "export",
			fmt.Sprintf("--model_type=%s", *task.ModelType),
			fmt.Sprintf("--model_id_or_path=%s", *task.ModelPath),
			fmt.Sprintf("--ckpt_dir=%s", *task.OutputDir),
			"--merge_lora=true",
			fmt.Sprintf("--output_dir=%s", outputPath),
			"--to_awq=true",
			fmt.Sprintf("--quant_bits=%d", quantizeBit),
		}
	case "gptq":
		exportCmd = []string{
			"swift", "export",
			fmt.Sprintf("--model_type=%s", *task.ModelType),
			fmt.Sprintf("--model_id_or_path=%s", *task.ModelPath),
			fmt.Sprintf("--ckpt_dir=%s", *task.OutputDir),
			"--merge_lora=true",
			fmt.Sprintf("--output_dir=%s", outputPath),
			"--to_gptq=true",
			fmt.Sprintf("--quant_bits=%d", quantizeBit),
		}
	default:
		return fmt.Errorf("不支持的导出格式: %s", exportFormat)
	}

	// 创建导出容器
	config := &container.Config{
		Image: "modelscope-registry.cn-hangzhou.cr.aliyuncs.com/modelscope-repo/modelscope:ubuntu22.04-cuda12.9.1-py311-torch2.8.0-vllm0.11.0-modelscope1.32.0-swift3.11.3",
		Cmd:   exportCmd,
		Tty:   false,
		Env: []string{
			"MODELSCOPE_CACHE=/data/models",
			"SWIFT_LOG_LEVEL=INFO",
		},
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "on-failure", MaximumRetryCount: 1},
		ShmSize:       8 * 1024 * 1024 * 1024,
	}

	// 挂载输出目录
	var binds []string
	outputBind := fmt.Sprintf("finetune-output-%d:%s", task.ID, *task.OutputDir)
	binds = append(binds, outputBind)
	hostConfig.Binds = binds

	containerName := fmt.Sprintf("finetune-export-%d-%d", task.ID, time.Now().Unix())

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return fmt.Errorf("创建导出容器失败: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("启动导出容器失败: %v", err)
	}

	global.GVA_LOG.Info("模型导出任务已启动", zap.Uint("taskId", taskID), zap.String("format", exportFormat))
	return nil
}

// CleanupOldTasks 清理7天前的已完成/失败任务
func (s *FineTuneTaskService) CleanupOldTasks(ctx context.Context) (count int64, err error) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Unix()

	result := global.GVA_DB.Where(
		"((task_status = ? AND end_time < ?) OR (task_status = ? AND end_time < ?)) AND deleted_at IS NULL",
		"completed", sevenDaysAgo,
		"failed", sevenDaysAgo,
	).Delete(&cloud.FineTuneTask{})

	count = result.RowsAffected
	global.GVA_LOG.Info("清理旧微调任务完成", zap.Int64("count", count))

	return count, nil
}

// SaveTaskSnapshotToFile 将快照数据保存到文件
func (s *FineTuneTaskService) SaveTaskSnapshotToFile(taskID uint) error {
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return err
	}

	snapshots, err := s.GetSnapshots(context.Background(), taskID)
	if err != nil {
		return err
	}

	// 创建快照目录
	snapshotDir := filepath.Join(*task.OutputDir, "snapshots")
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return err
	}

	// 保存快照数据
	snapshotFile := filepath.Join(snapshotDir, fmt.Sprintf("snapshots_%d.json", time.Now().Unix()))
	data, err := json.MarshalIndent(snapshots, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(snapshotFile, data, 0644)
}

// GetTaskStatistics 获取任务统计信息
func (s *FineTuneTaskService) GetTaskStatistics(ctx context.Context) (stats map[string]any, err error) {
	stats = make(map[string]any)

	// 按状态统计
	var statusStats []struct {
		TaskStatus string
		Count      int64
	}
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Select("task_status, count(*) as count").
		Where("deleted_at IS NULL").
		Group("task_status").
		Scan(&statusStats)

	statusMap := make(map[string]int64)
	for _, s := range statusStats {
		statusMap[s.TaskStatus] = s.Count
	}
	stats["byStatus"] = statusMap

	// 按模型类型统计
	var modelTypeStats []struct {
		ModelType string
		Count     int64
	}
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Select("model_type, count(*) as count").
		Where("deleted_at IS NULL").
		Group("model_type").
		Scan(&modelTypeStats)

	modelTypeMap := make(map[string]int64)
	for _, m := range modelTypeStats {
		modelTypeMap[m.ModelType] = m.Count
	}
	stats["byModelType"] = modelTypeMap

	// 今日任务统计
	todayStart := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Where("created_at >= ? AND deleted_at IS NULL", todayStart).
		Count(&todayCount)
	stats["todayCount"] = todayCount

	// 总任务数
	var totalCount int64
	global.GVA_DB.Model(&cloud.FineTuneTask{}).
		Where("deleted_at IS NULL").
		Count(&totalCount)
	stats["totalCount"] = totalCount

	return stats, nil
}

// CallSwiftAPI 调用SWIRT API (用于与SWIRT WebUI集成)
func (s *FineTuneTaskService) CallSwiftAPI(apiEndpoint string, payload map[string]any) (result map[string]any, err error) {
	// 默认SWIRT API地址
	swiftAPIURL := "http://127.0.0.1:17860" // 可从配置获取

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", swiftAPIURL+apiEndpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result = make(map[string]any)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API调用失败: %s", string(body))
	}

	return result, nil
}

// SyncSwiftTask 同步SWIRT WebUI的任务状态
func (s *FineTuneTaskService) SyncSwiftTask(ctx context.Context, taskID uint) error {
	var task cloud.FineTuneTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return err
	}

	// 调用SWIRT API获取任务状态
	payload := map[string]any{
		"task_id": taskID,
	}

	result, err := s.CallSwiftAPI("/api/task/status", payload)
	if err != nil {
		return err
	}

	// 解析返回的状态并更新数据库
	if status, ok := result["status"].(string); ok {
		task.TaskStatus = &status
	}
	if progress, ok := result["progress"].(float64); ok {
		task.Progress = &progress
	}
	if currentStep, ok := result["current_step"].(float64); ok {
		step := int64(currentStep)
		task.CurrentStep = &step
	}
	if totalSteps, ok := result["total_steps"].(float64); ok {
		steps := int64(totalSteps)
		task.TotalSteps = &steps
	}
	if loss, ok := result["loss"].(float64); ok {
		task.ValidationLoss = &loss
	}

	return global.GVA_DB.Save(&task).Error
}
