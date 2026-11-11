package docker_manager

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

// ansiPattern 终端颜色等控制序列，网页上显示没有意义
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[A-Za-z]`)

// MaxLogTail 单次最多返回的日志行数
const MaxLogTail = 5000

// splitLines 按 \n 或 \r 分行：steamcmd 的进度条用 \r 覆盖同一行
func splitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func cleanLine(line []byte) string {
	return ansiPattern.ReplaceAllString(string(line), "")
}

// readLines 把多路复用的日志流拆成行，空行跳过
func readLines(reader io.Reader, onLine func(string) error) error {
	pr, pw := io.Pipe()
	go func() {
		_, err := stdcopy.StdCopy(pw, pw, reader)
		pw.CloseWithError(err)
	}()
	defer pr.Close()

	scanner := bufio.NewScanner(pr)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	scanner.Split(splitLines)
	for scanner.Scan() {
		line := cleanLine(scanner.Bytes())
		if line == "" {
			continue
		}
		if err := onLine(line); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func clampTail(tail int) string {
	if tail <= 0 {
		tail = 200
	}
	if tail > MaxLogTail {
		tail = MaxLogTail
	}
	return strconv.Itoa(tail)
}

// ContainerLogTail 返回容器最近 tail 行日志（标准输出和标准错误合并）
func (dm *DockerManager) ContainerLogTail(containerName string, tail int) ([]string, error) {
	reader, err := dm.client.ContainerLogs(dm.ctx, containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       clampTail(tail),
	})
	if err != nil {
		return nil, fmt.Errorf("读取容器日志失败: %v", err)
	}
	defer reader.Close()

	lines := []string{}
	err = readLines(reader, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	return lines, err
}

// FollowContainerLogs 先输出最近 tail 行，然后持续输出新日志，直到 ctx 取消或容器停止
func (dm *DockerManager) FollowContainerLogs(ctx context.Context, containerName string, tail int, onLine func(string) error) error {
	reader, err := dm.client.ContainerLogs(ctx, containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       clampTail(tail),
	})
	if err != nil {
		return fmt.Errorf("读取容器日志失败: %v", err)
	}
	defer reader.Close()

	// ctx 取消时关闭读取流，让阻塞中的读取立即返回
	go func() {
		<-ctx.Done()
		reader.Close()
	}()

	err = readLines(reader, onLine)
	if ctx.Err() != nil {
		return nil
	}
	return err
}
