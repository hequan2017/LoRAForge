package cloud

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/volume"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"go.uber.org/zap"
)

type VolumeService struct{}

// GetVolumes 获取指定节点的卷列表
func (s *VolumeService) GetVolumes(ctx context.Context, nodeID int64) (volume.ListResponse, error) {
	// 1. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return volume.ListResponse{}, fmt.Errorf("未找到计算节点: %v", err)
	}

	// 2. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return volume.ListResponse{}, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 3. 获取卷列表
	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取卷列表失败", zap.Error(err))
		return volume.ListResponse{}, err
	}

	return volumes, nil
}

// CreateVolume 创建卷
func (s *VolumeService) CreateVolume(ctx context.Context, nodeID int64, name string, driver string) (volume.Volume, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return volume.Volume{}, fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	if driver == "" {
		driver = "local"
	}

	return cli.VolumeCreate(ctx, volume.CreateOptions{
		Name:   name,
		Driver: driver,
	})
}

// RemoveVolume 删除卷
func (s *VolumeService) RemoveVolume(ctx context.Context, nodeID int64, volumeID string, force bool) error {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	return cli.VolumeRemove(ctx, volumeID, force)
}
