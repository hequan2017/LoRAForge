package cloud

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"go.uber.org/zap"
)

type ImageService struct{}

// GetImages 获取指定节点的镜像列表
func (s *ImageService) GetImages(ctx context.Context, nodeID int64) ([]image.Summary, error) {
	// 1. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	// 2. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 3. 获取镜像列表
	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		global.GVA_LOG.Error("获取镜像列表失败", zap.Error(err))
		return nil, err
	}

	return images, nil
}

// PullImage 拉取镜像
func (s *ImageService) PullImage(ctx context.Context, nodeID int64, imageName string) (io.ReadCloser, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 自动补全 tag
	if !strings.Contains(imageName, ":") {
		imageName += ":latest"
	}

	return cli.ImagePull(ctx, imageName, image.PullOptions{})
}

// RemoveImage 删除镜像
func (s *ImageService) RemoveImage(ctx context.Context, nodeID int64, imageID string, force bool) ([]image.DeleteResponse, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	options := image.RemoveOptions{
		Force:         force,
		PruneChildren: true,
	}

	return cli.ImageRemove(ctx, imageID, options)
}
