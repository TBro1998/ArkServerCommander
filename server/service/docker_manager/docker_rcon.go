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

// rconScript 在 ARK 容器内通过 Source RCON 协议向 127.0.0.1 上的 RCON 端口发送一条命令。
// 后端容器和 ARK 容器不一定在同一个 Docker 网络里，所以不从后端直连，而是用 docker exec 在容器内部执行；
// 镜像已安装 python3，只用标准库。参数通过环境变量传入，密码不会出现在进程参数里。
const rconScript = `import os, socket, struct, sys

def packet(pid, ptype, body):
    data = struct.pack('<ii', pid, ptype) + body.encode() + b'\x00\x00'
    return struct.pack('<i', len(data)) + data

def recv_exact(s, n):
    buf = b''
    while len(buf) < n:
        chunk = s.recv(n - len(buf))
        if not chunk:
            raise IOError('connection closed')
        buf += chunk
    return buf

def read_packet(s):
    size = struct.unpack('<i', recv_exact(s, 4))[0]
    data = recv_exact(s, size)
    pid, ptype = struct.unpack('<ii', data[:8])
    return pid, ptype, data[8:-2].decode('utf-8', 'replace')

s = socket.create_connection(('127.0.0.1', int(os.environ['RCON_PORT'])), 10)
s.settimeout(float(os.environ.get('RCON_TIMEOUT', '60')))
s.sendall(packet(1, 3, os.environ['RCON_PASSWORD']))
while True:
    pid, ptype, body = read_packet(s)
    if ptype == 2:
        break
if pid == -1:
    sys.exit('RCON auth failed')
s.sendall(packet(2, 2, os.environ['RCON_COMMAND']))
pid, ptype, body = read_packet(s)
print(body)
`

// RconCommand 在容器内通过 RCON 执行一条命令，返回命令输出。
// timeout 是等待命令返回的时间（SaveWorld 在大地图上可能要几十秒）
func (dm *DockerManager) RconCommand(containerName string, port int, password, command string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(dm.ctx, timeout+15*time.Second)
	defer cancel()

	execResp, err := dm.client.ContainerExecCreate(ctx, containerName, container.ExecOptions{
		Cmd: []string{"python3", "-c", rconScript},
		Env: []string{
			fmt.Sprintf("RCON_PORT=%d", port),
			"RCON_PASSWORD=" + password,
			"RCON_COMMAND=" + command,
			fmt.Sprintf("RCON_TIMEOUT=%d", int(timeout.Seconds())),
		},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("创建RCON执行实例失败: %v", err)
	}

	attach, err := dm.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("执行RCON命令失败: %v", err)
	}
	defer attach.Close()

	// 非 TTY 的 exec 输出是多路复用的，需要拆分 stdout 和 stderr
	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, attach.Reader); err != nil {
		return "", fmt.Errorf("读取RCON输出失败: %v", err)
	}

	inspect, err := dm.client.ContainerExecInspect(ctx, execResp.ID)
	if err != nil {
		return "", fmt.Errorf("检查RCON执行结果失败: %v", err)
	}
	if inspect.ExitCode != 0 {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = fmt.Sprintf("退出码 %d", inspect.ExitCode)
		}
		return "", fmt.Errorf("RCON命令失败: %s", msg)
	}

	return strings.TrimSpace(stdout.String()), nil
}
