package cloud

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/network"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"go.uber.org/zap"
)

type NetworkService struct{}

// GetNetworks 获取指定节点的网络列表
func (s *NetworkService) GetNetworks(ctx context.Context, nodeID int64) ([]network.Inspect, error) {
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

	// 3. 获取网络列表
	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		global.GVA_LOG.Error("获取网络列表失败", zap.Error(err))
		return nil, err
	}

	return networks, nil
}

// CreateNetwork 创建网络
func (s *NetworkService) CreateNetwork(ctx context.Context, nodeID int64, name string, driver string) (string, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return "", fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return "", fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	if driver == "" {
		driver = "bridge"
	}

	resp, err := cli.NetworkCreate(ctx, name, network.CreateOptions{
		Driver: driver,
	})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// RemoveNetwork 删除网络
func (s *NetworkService) RemoveNetwork(ctx context.Context, nodeID int64, networkID string) error {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	return cli.NetworkRemove(ctx, networkID)
}
