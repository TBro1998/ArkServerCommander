package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ArkApi 和插件的日志写在服务端文件卷的 ShooterGame/Binaries/Win64/logs/*.log，
// 后端直接读文件卷，服务器停止后也能看

// maxPluginLogFiles 列出的日志文件数量上限（按修改时间从新到旧）
const maxPluginLogFiles = 200

// PluginLogFile 一个日志文件
type PluginLogFile struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}

// PluginLogFiles 日志文件列表
type PluginLogFiles struct {
	Files []PluginLogFile `json:"files"` // 从新到旧，最多 maxPluginLogFiles 个
	Total int             `json:"total"` // 目录里 .log 文件的总数
}

// pluginLogsDir 插件日志目录
func pluginLogsDir(serverID uint) string {
	return filepath.Join(arkApiWin64Dir(serverID), "logs")
}

// validPluginLogName 只能是 logs 目录下的单个 .log 文件名
func validPluginLogName(name string) bool {
	return name != "" && len(name) <= 255 && name == filepath.Base(name) && !strings.HasPrefix(name, ".") &&
		!strings.ContainsAny(name, `/\`) && strings.HasSuffix(strings.ToLower(name), ".log")
}

// ListPluginLogs 列出插件日志目录里的 .log 文件，目录不存在时返回空列表
func (s *ServerService) ListPluginLogs(userID uint, serverID string) (*PluginLogFiles, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	result := &PluginLogFiles{Files: []PluginLogFile{}}
	entries, err := os.ReadDir(pluginLogsDir(server.ID))
	if errors.Is(err, os.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取插件日志目录失败: %w", err)
	}
	type file struct {
		PluginLogFile
		mod time.Time
	}
	var files []file
	for _, e := range entries {
		if !e.Type().IsRegular() || !validPluginLogName(e.Name()) {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, file{PluginLogFile{Name: e.Name(), Size: fi.Size(), ModifiedAt: fi.ModTime().Format(time.RFC3339)}, fi.ModTime()})
	}
	sort.Slice(files, func(i, j int) bool {
		if !files[i].mod.Equal(files[j].mod) {
			return files[i].mod.After(files[j].mod)
		}
		return files[i].Name > files[j].Name
	})
	result.Total = len(files)
	if len(files) > maxPluginLogFiles {
		files = files[:maxPluginLogFiles]
	}
	for _, f := range files {
		result.Files = append(result.Files, f.PluginLogFile)
	}
	return result, nil
}

// GetPluginLog 返回一个插件日志文件的最后 tail 行（编码自动识别，同 ShooterGame.log）
func (s *ServerService) GetPluginLog(userID uint, serverID, name string, tail int) (*GameLog, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	if !validPluginLogName(name) {
		return nil, ErrPluginLogNotFound
	}
	// 目录内容由 ARK 容器写入，用 os.Root 打开，符号链接不能指到目录外
	root, err := os.OpenRoot(pluginLogsDir(server.ID))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrPluginLogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("读取插件日志失败: %w", err)
	}
	defer root.Close()
	f, err := root.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrPluginLogNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("读取插件日志失败: %w", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("读取插件日志失败: %w", err)
	}
	if !fi.Mode().IsRegular() {
		return nil, ErrPluginLogNotFound
	}
	lines, err := tailLogFile(f, fi.Size(), tail)
	if err != nil {
		return nil, fmt.Errorf("读取插件日志失败: %w", err)
	}
	return &GameLog{Exists: true, Lines: lines, Size: fi.Size(), ModifiedAt: fi.ModTime().Format(time.RFC3339)}, nil
}
