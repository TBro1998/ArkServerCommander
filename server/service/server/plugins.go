package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
)

// ArkApi 插件管理。插件卷绑定在数据目录 servers/<id>/plugins，后端直接读写；
// 停用的插件移到不挂载进容器的 servers/<id>/plugins-disabled，ArkApi 不会加载

const (
	maxPluginZipEntries = 5000
	maxPluginUnpacked   = 500 << 20 // 解压后总大小上限
	maxPluginConfigSize = 4 << 20   // config.json 大小上限
)

// pluginNamePattern 插件目录名（也是 dll 名称）
var pluginNamePattern = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_. -]{0,99}$`)

// pluginLocks 同一服务器的插件安装、删除、启停串行执行
var pluginLocks sync.Map // uint -> *sync.Mutex

func lockPlugins(serverID uint) func() {
	v, _ := pluginLocks.LoadOrStore(serverID, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// PluginInfo 一个已安装的插件
type PluginInfo struct {
	Name          string   `json:"name"`    // 目录名
	Enabled       bool     `json:"enabled"` // false 表示已停用（不在插件卷里）
	FullName      string   `json:"full_name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Version       string   `json:"version,omitempty"`
	MinAPIVersion string   `json:"min_api_version,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	HasDLL        bool     `json:"has_dll"` // 目录里有 <名称>.dll，没有时 ArkApi 加载会失败
	HasConfig     bool     `json:"has_config"`
	Size          int64    `json:"size"`
}

// PluginsOverview 插件列表和 ArkApi 本体信息
type PluginsOverview struct {
	ArkApiInstalled bool         `json:"arkapi_installed"`         // 服务端文件里有 ArkApi（version.dll）
	ArkApiVersion   string       `json:"arkapi_version,omitempty"` // 安装时写入的版本号（.arkapi_version），手动放进去的没有
	Plugins         []PluginInfo `json:"plugins"`
}

// ListPlugins 列出已安装的插件（含已停用的）和 ArkApi 版本
func (s *ServerService) ListPlugins(userID uint, serverID string) (*PluginsOverview, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	overview := &PluginsOverview{Plugins: []PluginInfo{}}
	overview.ArkApiVersion, overview.ArkApiInstalled = installedArkApiVersion(server.ID)

	for _, dir := range []struct {
		path    string
		enabled bool
	}{
		{docker_manager.ServerPluginsDir(server.ID), true},
		{docker_manager.ServerDisabledPluginsDir(server.ID), false},
	} {
		entries, err := os.ReadDir(dir.path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("读取插件目录失败: %w", err)
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			overview.Plugins = append(overview.Plugins, readPluginInfo(filepath.Join(dir.path, e.Name()), e.Name(), dir.enabled))
		}
	}
	sort.Slice(overview.Plugins, func(i, j int) bool {
		return strings.ToLower(overview.Plugins[i].Name) < strings.ToLower(overview.Plugins[j].Name)
	})
	return overview, nil
}

// readPluginInfo 读取插件目录里的 PluginInfo.json 和文件信息，读不到的字段留空
func readPluginInfo(dir, name string, enabled bool) PluginInfo {
	info := PluginInfo{Name: name, Enabled: enabled}
	if _, err := os.Stat(filepath.Join(dir, name+".dll")); err == nil {
		info.HasDLL = true
	}
	if _, err := os.Stat(filepath.Join(dir, "config.json")); err == nil {
		info.HasConfig = true
	}
	filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if fi, err := d.Info(); err == nil {
				info.Size += fi.Size()
			}
		}
		return nil
	})

	data, err := os.ReadFile(filepath.Join(dir, "PluginInfo.json"))
	if err != nil {
		return info
	}
	var meta struct {
		FullName      string   `json:"FullName"`
		Description   string   `json:"Description"`
		Version       any      `json:"Version"`
		MinAPIVersion any      `json:"MinApiVersion"`
		Dependencies  []string `json:"Dependencies"`
	}
	if json.Unmarshal(stripBOM(data), &meta) != nil {
		return info
	}
	info.FullName = meta.FullName
	info.Description = meta.Description
	info.Version = jsonScalar(meta.Version)
	info.MinAPIVersion = jsonScalar(meta.MinAPIVersion)
	info.Dependencies = meta.Dependencies
	return info
}

// jsonScalar 版本号可能写成数字或字符串
func jsonScalar(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return fmt.Sprint(x)
	}
	return ""
}

func stripBOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}

// findPlugin 找到插件所在的目录，返回路径和是否启用
func findPlugin(serverID uint, name string) (string, bool, error) {
	if !pluginNamePattern.MatchString(name) {
		return "", false, ErrPluginNotFound
	}
	enabled := filepath.Join(docker_manager.ServerPluginsDir(serverID), name)
	if fi, err := os.Stat(enabled); err == nil && fi.IsDir() {
		return enabled, true, nil
	}
	disabled := filepath.Join(docker_manager.ServerDisabledPluginsDir(serverID), name)
	if fi, err := os.Stat(disabled); err == nil && fi.IsDir() {
		return disabled, false, nil
	}
	return "", false, ErrPluginNotFound
}

// stoppedServer 插件文件的增删和启停只能在服务器停止时进行：运行中的游戏进程加载着这些 dll，替换会导致崩溃
func stoppedServer(userID uint, serverID string) (models.Server, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return server, err
	}
	if server.Status != "stopped" {
		return server, ErrServerNotStopped
	}
	return server, nil
}

// InstallPlugins 从 zip 安装插件，返回安装的插件名称。
// 压缩包里每个 <名称>/<名称>.dll 所在的目录是一个插件（可以有多层前缀，例如 ArkApi/Plugins/<名称>/）；
// 找不到这种结构时，name 指定插件名称，或者压缩包里只有一个 dll 时以它的名称为准。
// 覆盖已安装的插件时保留原来的 config.json 和启用状态
func (s *ServerService) InstallPlugins(userID uint, serverID string, r io.ReaderAt, size int64, name string) ([]string, error) {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("%w: 不是有效的 zip 文件", ErrInvalidPlugin)
	}
	roots, err := pluginRoots(zr.File, name)
	if err != nil {
		return nil, err
	}

	unlock := lockPlugins(server.ID)
	defer unlock()

	var installed []string
	for _, root := range roots {
		if err := installPlugin(server.ID, zr.File, root.prefix, root.name); err != nil {
			return installed, err
		}
		installed = append(installed, root.name)
	}
	return installed, nil
}

type pluginRoot struct {
	prefix string // 压缩包内的目录前缀，以 / 结尾，根目录为空
	name   string
}

// pluginRoots 识别压缩包里的插件目录，同时检查条目数量、大小和路径
func pluginRoots(files []*zip.File, name string) ([]pluginRoot, error) {
	if len(files) > maxPluginZipEntries {
		return nil, fmt.Errorf("%w: 压缩包文件太多", ErrInvalidPlugin)
	}
	var total uint64
	var dlls []string
	for _, f := range files {
		p := f.Name
		if strings.HasPrefix(p, "__MACOSX/") {
			continue
		}
		if !safeZipPath(p) {
			return nil, fmt.Errorf("%w: 压缩包包含不安全的路径 %s", ErrInvalidPlugin, p)
		}
		if f.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("%w: 压缩包不能包含符号链接", ErrInvalidPlugin)
		}
		total += f.UncompressedSize64
		if total > maxPluginUnpacked {
			return nil, fmt.Errorf("%w: 解压后超过 500 MB", ErrInvalidPlugin)
		}
		if !f.FileInfo().IsDir() && strings.EqualFold(path.Ext(p), ".dll") {
			dlls = append(dlls, p)
		}
	}

	dllName := func(p string) string { return strings.TrimSuffix(path.Base(p), path.Ext(p)) }
	dirOf := func(p string) string {
		if d := path.Dir(p); d != "." {
			return d + "/"
		}
		return ""
	}

	if name != "" {
		if !pluginNamePattern.MatchString(name) {
			return nil, fmt.Errorf("%w: 插件名称只能包含字母、数字、下划线、点、空格和连字符", ErrInvalidPlugin)
		}
		for _, p := range dlls {
			if strings.EqualFold(dllName(p), name) {
				return []pluginRoot{{prefix: dirOf(p), name: name}}, nil
			}
		}
		return nil, fmt.Errorf("%w: 压缩包里没有 %s.dll", ErrInvalidPlugin, name)
	}

	// <名称>/<名称>.dll：dll 与所在目录同名
	var roots []pluginRoot
	for _, p := range dlls {
		dir := path.Dir(p)
		if dir != "." && strings.EqualFold(path.Base(dir), dllName(p)) && pluginNamePattern.MatchString(path.Base(dir)) {
			roots = append(roots, pluginRoot{prefix: dir + "/", name: path.Base(dir)})
		}
	}
	// 嵌套在另一个插件目录里的（插件自带的子目录）不单独安装
	sort.Slice(roots, func(i, j int) bool { return len(roots[i].prefix) < len(roots[j].prefix) })
	var result []pluginRoot
	for _, r := range roots {
		nested := false
		for _, kept := range result {
			if strings.HasPrefix(r.prefix, kept.prefix) {
				nested = true
				break
			}
		}
		if !nested {
			result = append(result, r)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	if len(dlls) == 1 && pluginNamePattern.MatchString(dllName(dlls[0])) {
		return []pluginRoot{{prefix: dirOf(dlls[0]), name: dllName(dlls[0])}}, nil
	}
	return nil, fmt.Errorf("%w: 无法识别插件目录，压缩包里应包含 <插件名>/<插件名>.dll，或填写插件名称", ErrInvalidPlugin)
}

// safeZipPath 只允许相对路径，不能包含 .. 或盘符
func safeZipPath(p string) bool {
	if p == "" || strings.Contains(p, "\\") || strings.HasPrefix(p, "/") || strings.Contains(p, ":") {
		return false
	}
	for _, part := range strings.Split(p, "/") {
		if part == ".." {
			return false
		}
	}
	return true
}

// installPlugin 把压缩包里 prefix 下的文件解压到临时目录，再替换插件目录
func installPlugin(serverID uint, files []*zip.File, prefix, name string) error {
	base := docker_manager.ServerDir(serverID)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	// 临时目录放在插件卷外面，安装中断也不会被 ArkApi 当成插件加载
	tmp, err := os.MkdirTemp(base, ".plugin-install-")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tmp)

	for _, f := range files {
		if !strings.HasPrefix(f.Name, prefix) || strings.HasPrefix(f.Name, "__MACOSX/") {
			continue
		}
		rel := strings.TrimPrefix(f.Name, prefix)
		if rel == "" {
			continue
		}
		target := filepath.Join(tmp, filepath.FromSlash(rel))
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("解压插件失败: %w", err)
			}
			continue
		}
		if err := extractZipFile(f, target); err != nil {
			return err
		}
	}

	// MkdirTemp 建的目录是 0700，改成和其他插件目录一致
	if err := os.Chmod(tmp, 0o755); err != nil {
		return fmt.Errorf("安装插件失败: %w", err)
	}

	dest := filepath.Join(docker_manager.ServerPluginsDir(serverID), name)
	old, _, err := findPlugin(serverID, name)
	if err != nil {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("创建插件目录失败: %w", err)
		}
		if err := os.Rename(tmp, dest); err != nil {
			return fmt.Errorf("安装插件失败: %w", err)
		}
		return nil
	}

	// 覆盖安装：保留原来的配置和位置（启用状态），旧目录先移开，新目录放好后再删除
	if data, err := os.ReadFile(filepath.Join(old, "config.json")); err == nil {
		if err := os.WriteFile(filepath.Join(tmp, "config.json"), data, 0o644); err != nil {
			return fmt.Errorf("保留插件配置失败: %w", err)
		}
	}
	backup := tmp + "-old"
	if err := os.Rename(old, backup); err != nil {
		return fmt.Errorf("替换旧版本插件失败: %w", err)
	}
	if err := os.Rename(tmp, old); err != nil {
		os.Rename(backup, old)
		return fmt.Errorf("安装插件失败: %w", err)
	}
	os.RemoveAll(backup)
	return nil
}

func extractZipFile(f *zip.File, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("解压插件失败: %w", err)
	}
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("%w: 读取 %s 失败", ErrInvalidPlugin, f.Name)
	}
	defer rc.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("解压插件失败: %w", err)
	}
	// 按声明的大小限制实际写入量，防止伪造大小的压缩炸弹
	n, err := io.Copy(out, io.LimitReader(rc, int64(f.UncompressedSize64)+1))
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("%w: 解压 %s 失败", ErrInvalidPlugin, f.Name)
	}
	if uint64(n) > f.UncompressedSize64 {
		return fmt.Errorf("%w: %s 的实际大小与声明不符", ErrInvalidPlugin, f.Name)
	}
	return nil
}

// DeletePlugin 删除插件（包括配置）
func (s *ServerService) DeletePlugin(userID uint, serverID, name string) error {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return err
	}
	unlock := lockPlugins(server.ID)
	defer unlock()
	dir, _, err := findPlugin(server.ID, name)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("删除插件失败: %w", err)
	}
	return nil
}

// SetPluginEnabled 启用或停用插件：在插件卷和 plugins-disabled 之间移动目录
func (s *ServerService) SetPluginEnabled(userID uint, serverID, name string, enabled bool) error {
	server, err := stoppedServer(userID, serverID)
	if err != nil {
		return err
	}
	unlock := lockPlugins(server.ID)
	defer unlock()
	dir, current, err := findPlugin(server.ID, name)
	if err != nil {
		return err
	}
	if current == enabled {
		return nil
	}
	target := filepath.Join(docker_manager.ServerPluginsDir(server.ID), name)
	if !enabled {
		target = filepath.Join(docker_manager.ServerDisabledPluginsDir(server.ID), name)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("创建插件目录失败: %w", err)
	}
	if err := os.Rename(dir, target); err != nil {
		return fmt.Errorf("移动插件目录失败: %w", err)
	}
	return nil
}

// arkApiConfigName 表示 ArkApi 本体的配置，不是插件
const arkApiConfigName = ""

// pluginConfigPath 插件或 ArkApi 本体（name 为空）的 config.json 路径
func pluginConfigPath(serverID uint, name string) (string, error) {
	if name == arkApiConfigName {
		return filepath.Join(arkApiWin64Dir(serverID), "config.json"), nil
	}
	dir, _, err := findPlugin(serverID, name)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// GetPluginConfig 读取插件（name 为空时为 ArkApi 本体）的 config.json，文件不存在时返回空字符串
func (s *ServerService) GetPluginConfig(userID uint, serverID, name string) (string, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return "", err
	}
	p, err := pluginConfigPath(server.ID, name)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		if name == arkApiConfigName {
			return "", ErrArkApiNotInstalled
		}
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("读取配置失败: %w", err)
	}
	return string(stripBOM(data)), nil
}

// SavePluginConfig 保存 config.json，内容必须是合法 JSON。运行中也可以保存，插件重新加载或服务器重启后生效
func (s *ServerService) SavePluginConfig(userID uint, serverID, name, content string) error {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	if len(content) > maxPluginConfigSize {
		return fmt.Errorf("%w: 配置文件超过 4 MB", ErrInvalidPluginConfig)
	}
	if !json.Valid([]byte(content)) {
		return fmt.Errorf("%w: 不是合法的 JSON", ErrInvalidPluginConfig)
	}
	p, err := pluginConfigPath(server.ID, name)
	if err != nil {
		return err
	}
	if name == arkApiConfigName {
		if _, installed := installedArkApiVersion(server.ID); !installed {
			return ErrArkApiNotInstalled
		}
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("保存配置失败: %w", err)
	}
	return nil
}
