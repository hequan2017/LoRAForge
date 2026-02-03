package cloud

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
)

// ListContainerFiles 列出容器内文件
func (instService *InstanceService) ListContainerFiles(ctx context.Context, nodeID int64, containerID string, dirPath string) ([]map[string]interface{}, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	if dirPath == "" {
		dirPath = "/"
	}

	// 使用 exec 执行 ls -al
	cmd := []string{"ls", "-al", "--full-time", dirPath}
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return nil, err
	}

	// 解析 ls 输出
	files := make([]map[string]interface{}, 0)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// 简单解析，实际可能需要更严谨的解析逻辑
		// drwxr-xr-x 1 root root 4096 2021-01-01 00:00:00.000000000 +0000 .
		// 忽略 header total
		if strings.HasPrefix(line, "total") || line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue
		}
		
		isDir := strings.HasPrefix(parts[0], "d")
		size := parts[4]
		name := strings.Join(parts[8:], " ")
		
		// 忽略 . 和 ..
		if name == "." || name == ".." {
			continue
		}

		files = append(files, map[string]interface{}{
			"name":  name,
			"isDir": isDir,
			"size":  size,
			"perm":  parts[0],
			"user":  parts[2],
			"group": parts[3],
			"date":  parts[5] + " " + parts[6],
		})
	}

	return files, nil
}

// DownloadContainerFile 下载容器文件
func (instService *InstanceService) DownloadContainerFile(ctx context.Context, nodeID int64, containerID string, filePath string) (io.ReadCloser, string, error) {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return nil, "", fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return nil, "", fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	
	// CopyFromContainer 返回的是 tar 流
	reader, stat, err := cli.CopyFromContainer(ctx, containerID, filePath)
	if err != nil {
		cli.Close()
		return nil, "", err
	}
	
	// 我们需要封装一个 ReadCloser 来在关闭流时关闭 cli
	return &DockerStatsReader{ // 复用之前的结构体，或者新建一个
		Body: reader,
		Cli:  cli,
	}, stat.Name, nil
}

// UploadContainerFile 上传文件到容器
func (instService *InstanceService) UploadContainerFile(ctx context.Context, nodeID int64, containerID string, destPath string, fileHeader io.Reader, fileName string) error {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 制作 tar 包
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	
	content, err := io.ReadAll(fileHeader)
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: fileName,
		Mode: 0644,
		Size: int64(len(content)),
	}
	
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tw.Write(content); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}

	// 上传
	err = cli.CopyToContainer(ctx, containerID, destPath, buf, container.CopyToContainerOptions{})
	return err
}

// CreateDirectory 在容器中创建目录
func (instService *InstanceService) CreateDirectory(ctx context.Context, nodeID int64, containerID string, dirPath string) error {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 使用 mkdir -p
	cmd := []string{"mkdir", "-p", dirPath}
	execConfig := container.ExecOptions{
		Cmd: cmd,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return err
	}

	if err := cli.ContainerExecStart(ctx, execIDResp.ID, container.ExecStartOptions{}); err != nil {
		return err
	}

	return nil
}

// RemoveContainerFile 删除容器文件
func (instService *InstanceService) RemoveContainerFile(ctx context.Context, nodeID int64, containerID string, filePath string) error {
	var node cloud.ComputeNode
	if err := global.GVA_DB.First(&node, nodeID).Error; err != nil {
		return fmt.Errorf("未找到计算节点: %v", err)
	}

	cli, err := CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 使用 rm -rf
	cmd := []string{"rm", "-rf", filePath}
	execConfig := container.ExecOptions{
		Cmd: cmd,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return err
	}

	if err := cli.ContainerExecStart(ctx, execIDResp.ID, container.ExecStartOptions{}); err != nil {
		return err
	}

	return nil
}
