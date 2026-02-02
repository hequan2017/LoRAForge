package cloud

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	cloudReq "github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type InstanceService struct{}

// CreateInstance 创建实例管理记录
func (instService *InstanceService) CreateInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	global.GVA_LOG.Info("开始创建实例请求", zap.Any("instance", inst))

	// 0. 验证必填字段
	if inst.MirrorID == nil || *inst.MirrorID == 0 {
		global.GVA_LOG.Error("镜像ID不能为空")
		return fmt.Errorf("镜像ID不能为空")
	}
	if inst.NodeID == nil || *inst.NodeID == 0 {
		global.GVA_LOG.Error("节点ID不能为空")
		return fmt.Errorf("节点ID不能为空")
	}
	if inst.InstanceName == nil || *inst.InstanceName == "" {
		global.GVA_LOG.Error("实例名称不能为空")
		return fmt.Errorf("实例名称不能为空")
	}

	// 1. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, inst.NodeID).Error; err != nil {
		global.GVA_LOG.Error("获取计算节点失败", zap.Error(err))
		return fmt.Errorf("未找到计算节点: %v", err)
	}
	global.GVA_LOG.Info("获取到计算节点", zap.String("nodeName", *node.Name))

	// 2. 获取镜像信息
	var mirror cloud.MirrorRepository
	if err := global.GVA_DB.First(&mirror, inst.MirrorID).Error; err != nil {
		global.GVA_LOG.Error("获取镜像信息失败", zap.Error(err))
		return fmt.Errorf("未找到镜像: %v", err)
	}
	global.GVA_LOG.Info("获取到镜像信息", zap.String("mirror", *mirror.Address))

	// 3. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		global.GVA_LOG.Error("创建Docker客户端失败", zap.Error(err))
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	global.GVA_LOG.Info("Docker客户端创建成功")
	defer cli.Close()

	// 4. 准备容器配置
	imageName := *mirror.Address

	global.GVA_LOG.Info("开始拉取镜像", zap.String("image", imageName))
	// 尝试拉取镜像
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err == nil {
		global.GVA_LOG.Info("镜像拉取成功")
		defer reader.Close()
		io.Copy(io.Discard, reader)
	} else {
		global.GVA_LOG.Warn("镜像拉取失败或已存在", zap.Error(err))
	}

	// 解析配置
	config := &container.Config{
		Image: imageName,
		Tty:   true, // 分配伪终端，类似 -t
	}
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}

	// 实例名称
	containerName := ""
	if inst.InstanceName != nil && *inst.InstanceName != "" {
		// 校验并清理容器名称，仅允许 [a-zA-Z0-9][a-zA-Z0-9_.-]
		rawName := *inst.InstanceName
		// 简单的清理逻辑：将非法字符替换为 -
		var safeNameBuilder strings.Builder
		for _, r := range rawName {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' || r == '-' {
				safeNameBuilder.WriteRune(r)
			} else {
				safeNameBuilder.WriteRune('-')
			}
		}

		safeName := safeNameBuilder.String()
		// 去除连续的横线
		for strings.Contains(safeName, "--") {
			safeName = strings.ReplaceAll(safeName, "--", "-")
		}
		// 去除首尾的特殊字符
		safeName = strings.Trim(safeName, "_.-")

		if safeName == "" {
			safeName = fmt.Sprintf("%d", time.Now().Unix())
		}

		// 强制添加前缀，确保首字符合法且避免纯数字问题
		containerName = "inst-" + safeName
	}

	// 启动命令
	if inst.Command != nil && *inst.Command != "" {
		// 简单按空格分割，复杂命令可能需要更复杂的解析
		// 这里假设用户输入的是完整的命令字符串
		config.Cmd = strings.Fields(*inst.Command)
	}

	// 环境变量
	if inst.EnvVars != nil && *inst.EnvVars != "" {
		envs := strings.Split(*inst.EnvVars, "\n")
		var validEnvs []string
		for _, env := range envs {
			env = strings.TrimSpace(env)
			if env != "" {
				validEnvs = append(validEnvs, env)
			}
		}
		config.Env = validEnvs
	}

	// 端口映射
	if inst.PortMapping != nil && *inst.PortMapping != "" {
		exposedPorts := nat.PortSet{}
		portBindings := nat.PortMap{}

		ports := strings.Split(*inst.PortMapping, "\n")
		for _, p := range ports {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			// 格式 host:container
			parts := strings.Split(p, ":")
			if len(parts) == 2 {
				hostPort := parts[0]
				containerPort := parts[1]

				// 默认为 TCP
				portKey := nat.Port(containerPort + "/tcp")
				exposedPorts[portKey] = struct{}{}

				portBindings[portKey] = []nat.PortBinding{
					{
						HostPort: hostPort,
					},
				}
			}
		}
		config.ExposedPorts = exposedPorts
		hostConfig.PortBindings = portBindings
	}

	// 挂载目录
	if inst.VolumeMounts != nil && *inst.VolumeMounts != "" {
		var binds []string
		volumes := strings.Split(*inst.VolumeMounts, "\n")
		for _, v := range volumes {
			v = strings.TrimSpace(v)
			if v != "" {
				binds = append(binds, v)
			}
		}
		hostConfig.Binds = binds
	}

	// 资源限制 (CPU, Memory, GPU)
	// 优先使用自定义配置，如果未设置则尝试使用模版
	var cpuCores float64
	var memoryMB int64
	var gpuCount int64

	// 获取自定义配置
	if inst.Cpu != nil {
		cpuCores = *inst.Cpu
	}
	if inst.Memory != nil {
		memoryMB = *inst.Memory
	}
	if inst.GpuCount != nil {
		gpuCount = *inst.GpuCount
	}

	// 如果自定义未设置，尝试从模版获取
	if inst.TemplateID != nil && *inst.TemplateID > 0 {
		var template cloud.ProductSpec
		if err := global.GVA_DB.First(&template, inst.TemplateID).Error; err == nil {
			if cpuCores == 0 && template.CPUCores != nil {
				cpuCores = float64(*template.CPUCores)
			}
			if memoryMB == 0 && template.MemoryGB != nil {
				memoryMB = int64(*template.MemoryGB) * 1024 // GB to MB
			}
			if gpuCount == 0 && template.GPUCount != nil {
				gpuCount = int64(*template.GPUCount)
			}
		}
	}

	// 应用资源限制
	if cpuCores > 0 {
		hostConfig.Resources.NanoCPUs = int64(cpuCores * 1e9)
	}
	if memoryMB > 0 {
		hostConfig.Resources.Memory = memoryMB * 1024 * 1024
	}

	// GPU 配置
	if gpuCount > 0 {
		hostConfig.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Count:        int(gpuCount),
				Capabilities: [][]string{{"gpu"}},
			},
		}
	} else if gpuCount == -1 { // 这里的逻辑可以根据需求调整，比如 -1 代表 all
		hostConfig.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Count:        -1,
				Capabilities: [][]string{{"gpu"}},
			},
		}
	}

	// 5. 创建容器
	global.GVA_LOG.Info("准备创建容器",
		zap.String("name", containerName),
		zap.Any("config", config),
		zap.Any("hostConfig", hostConfig),
	)
	resp, err := cli.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		global.GVA_LOG.Error("Docker API创建容器失败", zap.Error(err))
		return fmt.Errorf("Docker API 创建容器失败: %v", err)
	}
	global.GVA_LOG.Info("容器创建成功", zap.String("id", resp.ID))

	// 6. 启动容器
	global.GVA_LOG.Info("准备启动容器", zap.String("id", resp.ID))
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		global.GVA_LOG.Error("Docker API启动容器失败", zap.Error(err))
		return fmt.Errorf("Docker API 启动容器失败: %v", err)
	}
	global.GVA_LOG.Info("容器启动成功")

	// 7. 更新数据库记录
	shortID := resp.ID
	if len(shortID) > 12 {
		shortID = shortID[:12]
	}

	inst.DockerContainer = &shortID
	status := "Running"
	inst.ContainerStatus = &status

	err = global.GVA_DB.Create(inst).Error
	return err
}

// CreateDockerClient 创建 Docker 客户端 (公开函数，供外部调用)
func CreateDockerClient(node cloud.ComputeNode) (*client.Client, error) {
	host := ""
	if node.DockerConnectAddress != nil {
		host = *node.DockerConnectAddress
	}
	if host == "" {
		return nil, fmt.Errorf("Docker连接地址为空")
	}

	// 自动补全 tcp:// 前缀
	if !strings.Contains(host, "://") {
		host = "tcp://" + host
	}

	var httpClient *http.Client
	if node.UseTLS != nil && *node.UseTLS {
		if node.CACertificate == nil || node.ClientCertificate == nil || node.ClientKey == nil {
			return nil, fmt.Errorf("TLS证书配置不完整")
		}

		caCert := []byte(*node.CACertificate)
		cert := []byte(*node.ClientCertificate)
		key := []byte(*node.ClientKey)

		// Load CA
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		// Load client cert
		certificate, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %v", err)
		}

		tlsConfig := &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{certificate},
			// InsecureSkipVerify: true, // 根据需要开启
		}

		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	}

	opts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	}
	if httpClient != nil {
		opts = append(opts, client.WithHTTPClient(httpClient))
	}

	return client.NewClientWithOpts(opts...)
}

// SyncInstances 同步算力节点的容器信息
func (instService *InstanceService) SyncInstances(ctx context.Context, nodeID int64) (err error) {
	// 1. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 2. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 3. 获取容器列表
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("获取容器列表失败: %v", err)
	}

	// 4. 同步数据库
	// 获取该节点下已存在的实例记录
	var existingInsts []cloud.Instance
	global.GVA_DB.Where("node_id = ?", nodeID).Find(&existingInsts)

	existingMap := make(map[string]*cloud.Instance)
	for i := range existingInsts {
		if existingInsts[i].DockerContainer != nil {
			// 以前缀12位为key，或者完整ID为key
			id := *existingInsts[i].DockerContainer
			if len(id) > 12 {
				id = id[:12]
			}
			existingMap[id] = &existingInsts[i]
		}
	}

	for _, c := range containers {
		// 匹配逻辑
		cID := c.ID
		if len(cID) > 12 {
			cID = cID[:12]
		}

		status := c.State
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		if inst, ok := existingMap[cID]; ok {
			// 更新状态
			inst.ContainerStatus = &status
			// 也可以选择更新名称等其他信息
			// inst.InstanceName = &name
			global.GVA_DB.Save(inst)
		} else {
			// 创建新记录
			// 注意：这里缺少 MirrorID, UserID 等信息，只能作为“发现”的实例
			// 建议设置一个默认的用户或标记为系统发现

			// 尝试解析端口映射
			portMapping := ""
			for _, p := range c.Ports {
				portMapping += fmt.Sprintf("%s:%d\n", p.IP, p.PublicPort)
			}

			newInst := cloud.Instance{
				GVA_MODEL:       global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				NodeID:          &nodeID,
				DockerContainer: &cID,
				InstanceName:    &name,
				ContainerStatus: &status,
				PortMapping:     &portMapping,
				Remark:          utils.Pointer("系统自动同步"),
			}
			global.GVA_DB.Create(&newInst)
		}
	}
	return nil
}

// CloseInstance 关闭实例
func (instService *InstanceService) CloseInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	// 1. 获取实例详情以拿到 NodeID 和 ContainerID
	var fullInst cloud.Instance
	if err := global.GVA_DB.First(&fullInst, inst.ID).Error; err != nil {
		return fmt.Errorf("未找到实例: %v", err)
	}
	if fullInst.NodeID == nil {
		return fmt.Errorf("实例未关联节点")
	}
	if fullInst.DockerContainer == nil || *fullInst.DockerContainer == "" {
		return fmt.Errorf("实例未关联容器")
	}

	// 2. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, fullInst.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 3. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 4. 停止容器
	// Timeout set to nil means wait indefinitely (or default timeout) for stop
	timeout := 30 // seconds
	if err := cli.ContainerStop(ctx, *fullInst.DockerContainer, container.StopOptions{Timeout: &timeout}); err != nil {
		global.GVA_LOG.Error("Docker API停止容器失败", zap.Error(err))
		return fmt.Errorf("Docker API 停止容器失败: %v", err)
	}

	// 5. 更新数据库状态
	status := "Stopped"
	return global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Update("container_status", status).Error
}

// RestartInstance 重启实例
func (instService *InstanceService) RestartInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	// 1. 获取实例详情
	var fullInst cloud.Instance
	if err := global.GVA_DB.First(&fullInst, inst.ID).Error; err != nil {
		return fmt.Errorf("未找到实例: %v", err)
	}
	if fullInst.NodeID == nil {
		return fmt.Errorf("实例未关联节点")
	}
	if fullInst.DockerContainer == nil || *fullInst.DockerContainer == "" {
		return fmt.Errorf("实例未关联容器")
	}

	// 2. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, fullInst.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 3. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 4. 重启容器
	if err := cli.ContainerRestart(ctx, *fullInst.DockerContainer, container.StopOptions{}); err != nil {
		global.GVA_LOG.Error("Docker API重启容器失败", zap.Error(err))
		return fmt.Errorf("Docker API 重启容器失败: %v", err)
	}

	// 5. 更新数据库状态
	status := "Running"
	return global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Update("container_status", status).Error
}

// StartInstance 启动实例
func (instService *InstanceService) StartInstance(ctx context.Context, inst *cloud.Instance) (err error) {
	// 1. 获取实例详情
	var fullInst cloud.Instance
	if err := global.GVA_DB.First(&fullInst, inst.ID).Error; err != nil {
		return fmt.Errorf("未找到实例: %v", err)
	}
	if fullInst.NodeID == nil {
		return fmt.Errorf("实例未关联节点")
	}
	if fullInst.DockerContainer == nil || *fullInst.DockerContainer == "" {
		return fmt.Errorf("实例未关联容器")
	}

	// 2. 获取节点信息
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, fullInst.NodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	// 3. 创建 Docker Client
	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 4. 启动容器
	if err := cli.ContainerStart(ctx, *fullInst.DockerContainer, container.StartOptions{}); err != nil {
		global.GVA_LOG.Error("Docker API启动容器失败", zap.Error(err))
		return fmt.Errorf("Docker API 启动容器失败: %v", err)
	}

	// 5. 更新数据库状态
	status := "Running"
	return global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Update("container_status", status).Error
}

// DeleteInstance 删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstance(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cloud.Instance{}, "id = ?", ID).Error
	return err
}

// DeleteInstanceByIds 批量删除实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) DeleteInstanceByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cloud.Instance{}, "id in ?", IDs).Error
	return err
}

// UpdateInstance 更新实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) UpdateInstance(ctx context.Context, inst cloud.Instance) (err error) {
	err = global.GVA_DB.Model(&cloud.Instance{}).Where("id = ?", inst.ID).Updates(&inst).Error
	return err
}

// GetInstance 根据ID获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstance(ctx context.Context, ID string) (inst cloud.Instance, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inst).Error
	return
}

// GetInstanceInfoList 分页获取实例管理记录
// Author [yourname](https://github.com/yourname)
func (instService *InstanceService) GetInstanceInfoList(ctx context.Context, info cloudReq.InstanceSearch) (list []cloud.Instance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cloud.Instance{})
	var insts []cloud.Instance
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.InstanceName != "" {
		db = db.Where("instance_name LIKE ?", "%"+info.InstanceName+"%")
	}
	if info.MirrorId != nil {
		db = db.Where("mirror_id = ?", *info.MirrorId)
	}
	if info.TemplateId != nil {
		db = db.Where("template_id = ?", *info.TemplateId)
	}
	if info.NodeId != nil {
		db = db.Where("node_id = ?", *info.NodeId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&insts).Error
	return insts, total, err
}
func (instService *InstanceService) GetInstanceDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)

	mirrorId := make([]map[string]any, 0)

	global.GVA_DB.Table("mirror_repositories").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&mirrorId)
	res["mirrorId"] = mirrorId
	nodeId := make([]map[string]any, 0)

	global.GVA_DB.Table("compute_nodes").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&nodeId)
	res["nodeId"] = nodeId
	templateId := make([]map[string]any, 0)

	global.GVA_DB.Table("product_specs").Where("deleted_at IS NULL").Select("name as label,id as value").Scan(&templateId)
	res["templateId"] = templateId
	return
}
func (instService *InstanceService) GetInstancePublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
