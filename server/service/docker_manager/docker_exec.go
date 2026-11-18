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

// ExecResult 容器内命令的执行结果
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ExecCommand 在容器内执行命令（argv[0] 是程序），返回标准输出、标准错误和退出码。
// 退出码不为 0 不算错误，由调用方判断；err 只表示无法执行或读取结果。timeout 是整个执行的超时时间
func (dm *DockerManager) ExecCommand(containerName string, argv []string, env []string, timeout time.Duration) (ExecResult, error) {
	ctx, cancel := context.WithTimeout(dm.ctx, timeout)
	defer cancel()

	execResp, err := dm.client.ContainerExecCreate(ctx, containerName, container.ExecOptions{
		Cmd:          argv,
		Env:          env,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return ExecResult{}, fmt.Errorf("创建执行实例失败: %v", err)
	}

	attach, err := dm.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return ExecResult{}, fmt.Errorf("执行命令失败: %v", err)
	}
	defer attach.Close()

	// 非 TTY 的 exec 输出是多路复用的，需要拆分 stdout 和 stderr
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, attach.Reader); err != nil {
		return ExecResult{}, fmt.Errorf("读取输出失败: %v", err)
	}

	inspect, err := dm.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return ExecResult{}, fmt.Errorf("检查执行结果失败: %v", err)
	}
	return ExecResult{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: inspect.ExitCode}, nil
}

// execPython 在 ARK 容器内用 python3 执行一段脚本，返回标准输出。
// 镜像已安装 python3；参数通过环境变量传入，不出现在进程参数里。
// timeout 是整个执行的超时时间
func (dm *DockerManager) execPython(containerName, script string, env []string, timeout time.Duration) (string, error) {
	res, err := dm.ExecCommand(containerName, []string{"python3", "-c", script}, env, timeout)
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		msg := strings.TrimSpace(res.Stderr)
		if msg == "" {
			msg = fmt.Sprintf("退出码 %d", res.ExitCode)
		}
		// 只保留最后一行，Python 的 traceback 对用户没有意义
		if idx := strings.LastIndex(msg, "\n"); idx >= 0 {
			msg = strings.TrimSpace(msg[idx+1:])
		}
		return "", fmt.Errorf("%s", msg)
	}

	return strings.TrimSpace(res.Stdout), nil
}
