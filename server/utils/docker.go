package utils

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// DockerManager Docker管理器结构体
type DockerManager struct {
	ctx context.Context
}

// NewDockerManager 创建新的Docker管理器
func NewDockerManager() *DockerManager {
	return &DockerManager{
		ctx: context.Background(),
	}
}

// CreateVolume 创建Docker卷
// serverID: 服务器ID
// 返回: 卷名称和错误信息
func (dm *DockerManager) CreateVolume(serverID uint) (string, error) {
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

	// 检查卷是否已存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return "", fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if exists {
		fmt.Printf("Docker卷 %s 已存在，跳过创建\n", volumeName)
		return volumeName, nil // 卷已存在，直接返回
	}

	// 创建卷
	fmt.Printf("正在创建Docker卷: %s\n", volumeName)
	cmd := exec.CommandContext(dm.ctx, "docker", "volume", "create", volumeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("创建Docker卷失败: %v, 输出: %s", err, string(output))
	}

	fmt.Printf("Docker卷创建成功: %s\n", volumeName)
	return volumeName, nil
}

// RemoveVolume 删除Docker卷
// volumeName: 卷名称
// 返回: 错误信息
func (dm *DockerManager) RemoveVolume(volumeName string) error {
	// 检查卷是否存在
	exists, err := dm.VolumeExists(volumeName)
	if err != nil {
		return fmt.Errorf("检查卷是否存在失败: %v", err)
	}
	if !exists {
		fmt.Printf("Docker卷 %s 不存在，跳过删除\n", volumeName)
		return nil // 卷不存在，视为删除成功
	}

	fmt.Printf("正在删除Docker卷: %s\n", volumeName)
	cmd := exec.CommandContext(dm.ctx, "docker", "volume", "rm", volumeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除Docker卷失败: %v, 输出: %s", err, string(output))
	}

	fmt.Printf("Docker卷删除成功: %s\n", volumeName)
	return nil
}

// VolumeExists 检查Docker卷是否存在
// volumeName: 卷名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) VolumeExists(volumeName string) (bool, error) {
	// 首先检查Docker是否可用
	testCmd := exec.CommandContext(dm.ctx, "docker", "version")
	if err := testCmd.Run(); err != nil {
		return false, fmt.Errorf("Docker不可用，请确保Docker已安装并运行: %v", err)
	}

	// 使用更简单的命令检查卷是否存在
	cmd := exec.CommandContext(dm.ctx, "docker", "volume", "inspect", volumeName)
	err := cmd.Run()
	if err != nil {
		// 如果卷不存在，docker volume inspect 会返回错误
		// 检查具体的错误信息
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil && strings.Contains(string(output), "No such volume") {
			return false, nil // 卷不存在，这是正常情况
		}
		// 其他错误情况
		return false, fmt.Errorf("检查Docker卷失败: %v, 输出: %s", err, string(output))
	}

	return true, nil
}

// CreateContainer 创建并启动ARK服务器容器
// serverID: 服务器ID
// serverName: 服务器名称
// port: 游戏端口
// queryPort: 查询端口
// rconPort: RCON端口
// adminPassword: 管理员密码
// mapName: 地图名称
// 返回: 容器ID和错误信息
func (dm *DockerManager) CreateContainer(serverID uint, serverName string, port, queryPort, rconPort int, adminPassword, mapName string) (string, error) {
	containerName := fmt.Sprintf("ark-server-%d", serverID)
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

	// 检查容器是否已存在
	if exists, err := dm.ContainerExists(containerName); err != nil {
		return "", fmt.Errorf("检查容器是否存在失败: %v", err)
	} else if exists {
		// 如果容器已存在，先删除它
		if err := dm.RemoveContainer(containerName); err != nil {
			return "", fmt.Errorf("删除已存在的容器失败: %v", err)
		}
	}

	// 构建docker run命令
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--restart", "unless-stopped",
		"-p", fmt.Sprintf("%d:7777/udp", port), // 游戏端口
		"-p", fmt.Sprintf("%d:7778/udp", port+1), // 原始套接字端口
		"-p", fmt.Sprintf("%d:27015/udp", queryPort), // 查询端口
		"-p", fmt.Sprintf("%d:32330/tcp", rconPort), // RCON端口
		"-v", fmt.Sprintf("%s:/ark", volumeName), // 挂载卷
		"-e", fmt.Sprintf("SESSIONNAME=%s", serverName),
		"-e", fmt.Sprintf("SERVERMAP=%s", mapName),
		"-e", fmt.Sprintf("SERVERPASSWORD="), // 服务器密码为空
		"-e", fmt.Sprintf("ADMINPASSWORD=%s", adminPassword),
		"-e", fmt.Sprintf("RCONENABLED=true"),
		"-e", fmt.Sprintf("RCONPORT=32330"), // 容器内部RCON端口固定为32330
		"-e", fmt.Sprintf("SERVERPORT=7777"), // 容器内部游戏端口固定为7777
		"-e", fmt.Sprintf("QUERYPORT=27015"), // 容器内部查询端口固定为27015
		"-e", "BACKUPONSTART=1", // 启动时备份
		"-e", "UPDATEONSTART=1", // 启动时更新
		"-e", "AUTOSTART=1", // 自动启动
		"tbro98/ase-server:v0.32", // 镜像名称
	}

	cmd := exec.CommandContext(dm.ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("创建Docker容器失败: %v, 输出: %s", err, string(output))
	}

	containerID := strings.TrimSpace(string(output))
	return containerID, nil
}

// StartContainer 启动容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StartContainer(containerName string) error {
	cmd := exec.CommandContext(dm.ctx, "docker", "start", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("启动Docker容器失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// StopContainer 停止容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) StopContainer(containerName string) error {
	// 使用优雅停止，给容器30秒时间保存数据
	cmd := exec.CommandContext(dm.ctx, "docker", "stop", "-t", "30", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止Docker容器失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// RemoveContainer 删除容器
// containerName: 容器名称
// 返回: 错误信息
func (dm *DockerManager) RemoveContainer(containerName string) error {
	// 先尝试停止容器
	dm.StopContainer(containerName)

	// 删除容器
	cmd := exec.CommandContext(dm.ctx, "docker", "rm", "-f", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除Docker容器失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// ContainerExists 检查容器是否存在
// containerName: 容器名称
// 返回: 是否存在和错误信息
func (dm *DockerManager) ContainerExists(containerName string) (bool, error) {
	// 使用inspect命令检查容器是否存在，这比filter更可靠
	cmd := exec.CommandContext(dm.ctx, "docker", "container", "inspect", containerName)
	err := cmd.Run()
	if err != nil {
		// 检查是否是"容器不存在"的错误
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil && strings.Contains(string(output), "No such container") {
			return false, nil // 容器不存在，这是正常情况
		}
		// 其他错误情况
		return false, fmt.Errorf("检查Docker容器失败: %v, 输出: %s", err, string(output))
	}

	return true, nil
}

// GetContainerStatus 获取容器状态
// containerName: 容器名称
// 返回: 状态字符串和错误信息
func (dm *DockerManager) GetContainerStatus(containerName string) (string, error) {
	cmd := exec.CommandContext(dm.ctx, "docker", "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", containerName), "--format", "{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取Docker容器状态失败: %v", err)
	}

	status := strings.TrimSpace(string(output))
	if status == "" {
		return "not_found", nil
	}

	// 解析Docker状态为我们的状态格式
	if strings.Contains(status, "Up") {
		return "running", nil
	} else if strings.Contains(status, "Exited") {
		return "stopped", nil
	} else if strings.Contains(status, "Created") {
		return "stopped", nil
	} else if strings.Contains(status, "Restarting") {
		return "starting", nil
	} else {
		return "unknown", nil
	}
}

// ExecuteCommand 在容器中执行命令
// containerName: 容器名称
// command: 要执行的命令
// 返回: 命令输出和错误信息
func (dm *DockerManager) ExecuteCommand(containerName string, command string) (string, error) {
	cmd := exec.CommandContext(dm.ctx, "docker", "exec", containerName, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("在容器中执行命令失败: %v, 输出: %s", err, string(output))
	}

	return string(output), nil
}

// ReadConfigFile 从容器卷中读取配置文件
// serverID: 服务器ID
// fileName: 文件名
// 返回: 文件内容和错误信息
func (dm *DockerManager) ReadConfigFile(serverID uint, fileName string) (string, error) {
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

	// 使用临时容器读取文件
	configPath := fmt.Sprintf("/ark/ShooterGame/Saved/Config/LinuxServer/%s", fileName)

	cmd := exec.CommandContext(dm.ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/ark", volumeName),
		"alpine:latest",
		"cat", configPath)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}

	return string(output), nil
}

// WriteConfigFile 向容器卷中写入配置文件
// serverID: 服务器ID
// fileName: 文件名
// content: 文件内容
// 返回: 错误信息
func (dm *DockerManager) WriteConfigFile(serverID uint, fileName, content string) error {
	volumeName := fmt.Sprintf("ark-server-%d", serverID)

	// 使用临时容器写入文件
	configPath := fmt.Sprintf("/ark/ShooterGame/Saved/Config/LinuxServer/%s", fileName)
	configDir := "/ark/ShooterGame/Saved/Config/LinuxServer"

	// 首先确保目录存在
	cmd := exec.CommandContext(dm.ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/ark", volumeName),
		"alpine:latest",
		"mkdir", "-p", configDir)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 写入文件内容
	cmd = exec.CommandContext(dm.ctx, "docker", "run", "--rm", "-i",
		"-v", fmt.Sprintf("%s:/ark", volumeName),
		"alpine:latest",
		"sh", "-c", fmt.Sprintf("cat > %s", configPath))

	cmd.Stdin = strings.NewReader(content)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// GetServerContainerName 获取服务器容器名称
// serverID: 服务器ID
// 返回: 容器名称
func GetServerContainerName(serverID uint) string {
	return fmt.Sprintf("ark-server-%d", serverID)
}

// GetServerVolumeName 获取服务器卷名称
// serverID: 服务器ID
// 返回: 卷名称
func GetServerVolumeName(serverID uint) string {
	return fmt.Sprintf("ark-server-%d", serverID)
}

// CheckDockerStatus 检查Docker环境状态
// 返回: Docker是否可用和错误信息
func CheckDockerStatus() error {
	// 检查Docker是否安装
	cmd := exec.Command("docker", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Docker未安装或不在PATH中: %v", err)
	}

	// 检查Docker服务是否运行
	cmd = exec.Command("docker", "info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Docker服务未运行: %v, 输出: %s", err, string(output))
	}

	return nil
}
