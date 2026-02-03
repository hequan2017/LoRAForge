package cloud

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"go.uber.org/zap"
)

type SwiftWebUIService struct{}

// CreateSwiftWebUI 创建Swift WebUI管理记录并启动容器
func (s *SwiftWebUIService) CreateSwiftWebUI(ctx context.Context, webui *cloud.SwiftWebUI) (err error) {
	// 1. 验证
	if webui.NodeID == nil {
		return fmt.Errorf("节点ID不能为空")
	}

	// 2. 获取节点
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, webui.NodeID).Error; err != nil {
		return fmt.Errorf("获取节点失败: %v", err)
	}

	// 3. 创建Docker客户端
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建Docker客户端失败: %v", err)
	}
	defer cli.Close()

	// 4. 镜像
	imageName := "modelscope-registry.cn-hangzhou.cr.aliyuncs.com/modelscope-repo/modelscope:ubuntu22.04-cuda12.9.1-py311-torch2.8.0-vllm0.11.0-modelscope1.32.0-swift3.11.3"
	// 尝试拉取镜像，忽略错误（可能已存在）
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err == nil {
		defer reader.Close()
		io.Copy(io.Discard, reader)
	}

	// 5. 准备启动命令
	// swift web-ui --lang zh --port 7860 --host 0.0.0.0
	lang := "zh"
	if webui.Language != nil {
		lang = *webui.Language
	}
	cmd := []string{"swift", "web-ui", "--lang", lang, "--port", "7860", "--host", "0.0.0.0"}

	// 6. 端口映射
	hostPort := "7860"
	if webui.Port != nil {
		hostPort = fmt.Sprintf("%d", *webui.Port)
	}

	portMap := nat.PortMap{
		"7860/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		},
	}

	// 7. 创建容器
	containerName := fmt.Sprintf("swift-webui-%d-%d", time.Now().Unix(), *webui.NodeID)

	hostConfig := &container.HostConfig{
		PortBindings: portMap,
		ShmSize:      8 * 1024 * 1024 * 1024, // 8GB
		Resources: container.Resources{
			DeviceRequests: []container.DeviceRequest{
				{
					Count:        -1, // 使用所有GPU
					Capabilities: [][]string{{"gpu"}},
				},
			},
		},
	}

	global.GVA_LOG.Info("Creating Swift WebUI container", zap.String("name", containerName), zap.Strings("cmd", cmd))

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Cmd:          cmd,
		ExposedPorts: nat.PortSet{"7860/tcp": struct{}{}},
		Tty:          true,
	}, hostConfig, &network.NetworkingConfig{}, nil, containerName)

	if err != nil {
		return fmt.Errorf("创建容器失败: %v", err)
	}

	// 8. 启动
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("启动容器失败: %v", err)
	}

	// 9. 更新信息
	status := "running"
	webui.Status = &status
	webui.ContainerId = &resp.ID

	ip := "localhost"
	if node.PublicIP != nil && *node.PublicIP != "" {
		ip = *node.PublicIP
	}
	// 如果是本机，可能需要调整IP logic，暂时使用PublicIP
	url := fmt.Sprintf("http://%s:%s", ip, hostPort)
	webui.AccessUrl = &url

	err = global.GVA_DB.Create(webui).Error
	return err
}

// DeleteSwiftWebUI 删除Swift WebUI管理记录
func (s *SwiftWebUIService) DeleteSwiftWebUI(ctx context.Context, ID string) (err error) {
	var webui cloud.SwiftWebUI
	if err := global.GVA_DB.First(&webui, ID).Error; err != nil {
		return err
	}

	// 停止容器
	if webui.ContainerId != nil && *webui.ContainerId != "" && webui.NodeID != nil {
		var node cloud.ComputeNode
		if err := global.GVA_DB.First(&node, webui.NodeID).Error; err == nil {
			if cli, err := CreateDockerClient(node); err == nil {
				defer cli.Close()
				// 停止并删除容器
				global.GVA_LOG.Info("Stopping container", zap.String("id", *webui.ContainerId))
				cli.ContainerRemove(ctx, *webui.ContainerId, container.RemoveOptions{Force: true})
			}
		}
	}

	err = global.GVA_DB.Delete(&webui).Error
	return err
}

// DeleteSwiftWebUIByIds 批量删除Swift WebUI管理记录
func (s *SwiftWebUIService) DeleteSwiftWebUIByIds(ctx context.Context, IDs []string) (err error) {
	for _, id := range IDs {
		_ = s.DeleteSwiftWebUI(ctx, id)
	}
	return nil
}

// UpdateSwiftWebUI 更新Swift WebUI管理记录
func (s *SwiftWebUIService) UpdateSwiftWebUI(ctx context.Context, webui cloud.SwiftWebUI) (err error) {
	err = global.GVA_DB.Model(&cloud.SwiftWebUI{}).Where("id = ?", webui.ID).Updates(&webui).Error
	return err
}

// GetSwiftWebUI 根据ID获取Swift WebUI管理记录
func (s *SwiftWebUIService) GetSwiftWebUI(ctx context.Context, ID string) (webui cloud.SwiftWebUI, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&webui).Error
	return
}

// GetSwiftWebUIInfoList 分页获取Swift WebUI管理记录
func (s *SwiftWebUIService) GetSwiftWebUIInfoList(ctx context.Context, info cloudReq.SwiftWebUISearch) (list []cloud.SwiftWebUI, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cloud.SwiftWebUI{})
	var webuis []cloud.SwiftWebUI
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.TaskName != nil && *info.TaskName != "" {
		db = db.Where("task_name LIKE ?", "%"+*info.TaskName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&webuis).Error
	return webuis, total, err
}

func (s *SwiftWebUIService) GetSwiftWebUIDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	nodeId := make([]map[string]any, 0)

	global.GVA_DB.Table("compute_nodes").Select("name as label,id as value").Scan(&nodeId)
	res["nodeId"] = nodeId
	return
}

func (s *SwiftWebUIService) GetSwiftWebUIPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
