package docker_manager

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// installedModsFile 入口脚本记录已安装 Mod 的文件，位于服务端文件卷根目录
const installedModsFile = ".installed_mods"

// parseInstalledMods 解析 .installed_mods：新格式每行 "id 安装时间（Unix 秒）"，旧格式只有 id（时间记为 0）
func parseInstalledMods(r io.Reader) map[string]int64 {
	result := make(map[string]int64)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		var ts int64
		if len(fields) > 1 {
			ts, _ = strconv.ParseInt(fields[1], 10, 64)
		}
		// 同一个 Mod 重新安装过时会有多行，以最后一次为准
		result[fields[0]] = ts
	}
	return result
}

// InstalledMods 读取服务器已安装的 Mod，返回 Mod ID -> 安装时间（Unix 秒）。
// 旧版入口脚本没有记录时间，用安装时写入的 Content/Mods/<id>.mod 的修改时间代替，仍然拿不到时为 0。
// 服务端文件还没有下载时返回空表
func InstalledMods(serverID uint) (map[string]int64, error) {
	dir := ServerFilesDir(serverID)
	f, err := os.Open(filepath.Join(dir, installedModsFile))
	if errors.Is(err, os.ErrNotExist) {
		return map[string]int64{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取已安装 Mod 记录失败: %w", err)
	}
	defer f.Close()

	mods := parseInstalledMods(f)
	for id, ts := range mods {
		if ts > 0 {
			continue
		}
		if info, err := os.Stat(filepath.Join(dir, "ShooterGame", "Content", "Mods", id+".mod")); err == nil {
			mods[id] = info.ModTime().Unix()
		}
	}
	return mods, nil
}

// RemoveInstalledMods 从 .installed_mods 删除指定的 Mod，下次启动时入口脚本只重新下载这些 Mod
func RemoveInstalledMods(serverID uint, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	path := filepath.Join(ServerFilesDir(serverID), installedModsFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("读取已安装 Mod 记录失败: %w", err)
	}

	remove := make(map[string]bool, len(ids))
	for _, id := range ids {
		remove[id] = true
	}
	var kept []string
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || remove[fields[0]] {
			continue
		}
		kept = append(kept, strings.TrimSpace(line))
	}
	content := strings.Join(kept, "\n")
	if content != "" {
		content += "\n"
	}

	// 先写临时文件再改名，避免写到一半时入口脚本读到不完整的内容
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return fmt.Errorf("写入已安装 Mod 记录失败: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("写入已安装 Mod 记录失败: %w", err)
	}
	return nil
}
