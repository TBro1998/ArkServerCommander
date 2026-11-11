package docker_manager

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

// execPython 在 ARK 容器内用 python3 执行一段脚本，返回标准输出。
// 镜像已安装 python3；参数通过环境变量传入，不出现在进程参数里。
// timeout 是整个执行的超时时间
func (dm *DockerManager) execPython(containerName, script string, env []string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(dm.ctx, timeout)
	defer cancel()

	execResp, err := dm.client.ContainerExecCreate(ctx, containerName, container.ExecOptions{
		Cmd:          []string{"python3", "-c", script},
		Env:          env,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("创建执行实例失败: %v", err)
	}

	attach, err := dm.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v", err)
	}
	defer attach.Close()

	// 非 TTY 的 exec 输出是多路复用的，需要拆分 stdout 和 stderr
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, attach.Reader); err != nil {
		return "", fmt.Errorf("读取输出失败: %v", err)
	}

	inspect, err := dm.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return "", fmt.Errorf("检查执行结果失败: %v", err)
	}
	if inspect.ExitCode != 0 {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = fmt.Sprintf("退出码 %d", inspect.ExitCode)
		}
		// 只保留最后一行，Python 的 traceback 对用户没有意义
		if idx := strings.LastIndex(msg, "\n"); idx >= 0 {
			msg = strings.TrimSpace(msg[idx+1:])
		}
		return "", fmt.Errorf("%s", msg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// RunHelper 用 alpine 临时容器执行一条命令（不经过 shell），等待结束后删除容器。
// binds 是卷挂载，如 "ase-server-files-1:/files"。返回命令的输出
func (dm *DockerManager) RunHelper(binds []string, cmd []string, timeout time.Duration) (string, error) {
	exists, err := dm.ImageExists(configHelperImage)
	if err != nil {
		return "", fmt.Errorf("检查Alpine镜像失败: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("Alpine镜像不存在，请先在首页下载镜像")
	}

	ctx, cancel := context.WithTimeout(dm.ctx, timeout)
	defer cancel()

	resp, err := dm.client.ContainerCreate(ctx,
		&container.Config{Image: configHelperImage, Cmd: cmd},
		&container.HostConfig{Binds: binds},
		nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("创建临时容器失败: %v", err)
	}
	defer func() {
		// 用独立的 context 清理，超时后也要删掉临时容器
		_ = dm.client.ContainerRemove(dm.ctx, resp.ID, container.RemoveOptions{Force: true})
	}()

	waitCh, errCh := dm.client.ContainerWait(ctx, resp.ID, container.WaitConditionNextExit)
	if err := dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("启动临时容器失败: %v", err)
	}

	var exitCode int64
	select {
	case result := <-waitCh:
		exitCode = result.StatusCode
	case err := <-errCh:
		return "", fmt.Errorf("等待临时容器失败: %v", err)
	}

	output := ""
	if logs, err := dm.client.ContainerLogs(dm.ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true}); err == nil {
		var stdout, stderr bytes.Buffer
		_, _ = stdcopy.StdCopy(&stdout, &stderr, logs)
		logs.Close()
		output = strings.TrimSpace(stdout.String() + stderr.String())
	}

	if exitCode != 0 {
		return output, fmt.Errorf("命令执行失败（退出码 %d）: %s", exitCode, output)
	}
	return output, nil
}
