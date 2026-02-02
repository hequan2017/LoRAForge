package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/image"
	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	cloudModel "github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudService "github.com/flipped-aurora/gin-vue-admin/server/service/cloud"
)

// SyncAllCloudStatus 同步所有算力节点的容器和镜像状态
// 定时任务调用此函数来保持云资源状态与实际 Docker 状态一致
func SyncAllCloudStatus() error {
	ctx := context.Background()
	startTime := time.Now()

	global.GVA_LOG.Info("开始同步云资源状态", zap.Time("startTime", startTime))

	// 1. 获取所有上架的算力节点
	var nodes []cloudModel.ComputeNode
	if err := global.GVA_DB.Where("is_listed = ?", true).Find(&nodes).Error; err != nil {
		global.GVA_LOG.Error("获取算力节点列表失败", zap.Error(err))
		return fmt.Errorf("获取算力节点列表失败: %v", err)
	}

	if len(nodes) == 0 {
		global.GVA_LOG.Info("没有找到上架的算力节点")
		return nil
	}

	global.GVA_LOG.Info("找到算力节点", zap.Int("count", len(nodes)))

	// 2. 遍历每个节点进行同步
	successCount := 0
	failCount := 0

	for _, node := range nodes {
		nodeName := "unknown"
		if node.Name != nil {
			nodeName = *node.Name
		}

		global.GVA_LOG.Debug("开始同步节点", zap.String("node", nodeName))

		// 同步容器状态
		if err := syncNodeContainers(ctx, node); err != nil {
			global.GVA_LOG.Error("同步节点容器失败",
				zap.String("node", nodeName),
				zap.Error(err))
			failCount++
			continue
		}

		// 同步镜像信息
		if err := syncNodeImages(ctx, node); err != nil {
			global.GVA_LOG.Error("同步节点镜像失败",
				zap.String("node", nodeName),
				zap.Error(err))
			// 镜像同步失败不影响整体结果，继续处理
		}

		successCount++
	}

	duration := time.Since(startTime)
	global.GVA_LOG.Info("云资源状态同步完成",
		zap.Int("total", len(nodes)),
		zap.Int("success", successCount),
		zap.Int("failed", failCount),
		zap.Duration("duration", duration))

	return nil
}

// syncNodeContainers 同步单个节点的容器状态
func syncNodeContainers(ctx context.Context, node cloudModel.ComputeNode) error {
	instanceService := &cloudService.InstanceService{}
	// node.ID 是 uint 类型，需要转换为 int64
	nodeID := int64(node.ID)
	return instanceService.SyncInstances(ctx, nodeID)
}

// syncNodeImages 同步单个节点的镜像信息到镜像库
func syncNodeImages(ctx context.Context, node cloudModel.ComputeNode) error {
	// 1. 创建 Docker Client
	cli, err := cloudService.CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 2. 获取镜像列表
	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("获取镜像列表失败: %v", err)
	}

	global.GVA_LOG.Debug("节点镜像列表",
		zap.String("node", safeString(node.Name)),
		zap.Int("count", len(images)))

	// 3. 同步到镜像库表
	// 获取该节点已存在的镜像记录
	var existingMirrors []cloudModel.MirrorRepository
	global.GVA_DB.Where("source = ?", fmt.Sprintf("node-%d", node.ID)).Find(&existingMirrors)

	existingMap := make(map[string]*cloudModel.MirrorRepository)
	for i := range existingMirrors {
		if existingMirrors[i].Address != nil {
			existingMap[*existingMirrors[i].Address] = &existingMirrors[i]
		}
	}

	// 遍历 Docker 镜像，同步或创建记录
	for _, img := range images {
		// 获取镜像标签
		imageTags := img.RepoTags
		if len(imageTags) == 0 {
			// 跳过没有标签的镜像（通常是中间层）
			continue
		}

		for _, tag := range imageTags {
			// 检查是否已存在（使用完整地址作为 key）
			if mirror, ok := existingMap[tag]; ok {
				// 已存在，更新时间戳
				global.GVA_DB.Model(mirror).Update("updated_at", time.Now())
			} else {
				// 创建新镜像记录
				// 提取镜像名称缩写：去掉仓库前缀和 tag
				name := shortenImageName(tag)
				desc := fmt.Sprintf("大小: %d MB", img.Size/1024/1024)

				newMirror := cloudModel.MirrorRepository{
					Name:        &name,
					Address:     &tag,
					Description: &desc,
					Source:      pointerString(fmt.Sprintf("node-%d", node.ID)),
					IsListed:    pointerBool(true), // 自动发现的镜像默认上架
					Remark:      pointerString(fmt.Sprintf("从节点 %s 自动同步", safeString(node.Name))),
				}
				if err := global.GVA_DB.Create(&newMirror).Error; err != nil {
					global.GVA_LOG.Error("创建镜像记录失败",
						zap.String("image", tag),
						zap.Error(err))
				} else {
					global.GVA_LOG.Info("自动添加镜像", zap.String("image", tag))
				}
			}
		}
	}

	return nil
}

// safeString 安全获取字符串值
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// pointerString 返回字符串指针的辅助函数
func pointerString(s string) *string {
	return &s
}

// pointerBool 返回布尔指针的辅助函数
func pointerBool(b bool) *bool {
	return &b
}

// shortenImageName 缩写镜像名称
// 例如：
// - "ubuntu:20.04" -> "ubuntu"
// - "nginx:latest" -> "nginx"
// - "docker.io/library/ubuntu:20.04" -> "ubuntu"
// - "registry.cn-hangzhou.aliyuncs.com/library/nginx:latest" -> "nginx"
func shortenImageName(fullName string) string {
	// 去掉 tag（:后面的部分）
	if idx := strings.Index(fullName, ":"); idx != -1 {
		fullName = fullName[:idx]
	}

	// 去掉仓库前缀，只保留最后一部分
	parts := strings.Split(fullName, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		if name != "" {
			return name
		}
	}

	return fullName
}
