package docker_manager

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ark-server-commander/utils"
)

// 服务端主程序相对安装目录的路径，存在才说明游戏本体已下载。Windows 版（Proton 镜像）和 Linux 版的主程序不同，
// 但在 Steam 上属于同一个 depot，下载时两个一起装上；任一个存在都算已下载
const (
	windowsServerExecutable = "ShooterGame/Binaries/Win64/ShooterGameServer.exe"
	linuxServerExecutable   = "ShooterGame/Binaries/Linux/ShooterGameServer"
)

// installMarkerFile 入口脚本下载完成后写入的安装标记，位于服务端文件卷根目录
const installMarkerFile = ".install_complete"

// gameBranchFile 文件卷里的服务端文件属于哪个游戏版本（Steam 分支），由管理端在启动前写入。
// 没有这个文件（加入游戏版本之前下载的）视为正式版
const gameBranchFile = ".game_branch"

// appManifestFile steamcmd 在安装目录里记录的安装信息，切换的 beta 分支（BetaKey）也记在里面
var appManifestFile = filepath.Join("steamapps", "appmanifest_"+utils.ArkServerAppID+".acf")

// betaKeyLine appmanifest 里记录 beta 分支的行
var betaKeyLine = regexp.MustCompile(`(?m)^[ \t]*"BetaKey"[ \t]+"[^"\n]*"[ \t]*\r?\n`)

// ServerFilesState 服务端文件卷的安装状态
type ServerFilesState struct {
	Installed bool   // 有 .install_complete 标记（入口脚本下载完成后写入）
	Windows   bool   // Windows 版服务端主程序存在
	Linux     bool   // Linux 版服务端主程序存在
	Branch    string // 服务端文件的游戏版本（.game_branch），空为正式版
}

// HasExecutable 任一平台的服务端主程序存在，说明游戏本体已经下载过
func (s ServerFilesState) HasExecutable() bool {
	return s.Windows || s.Linux
}

// isRegularFile 路径存在且是普通文件（跟随符号链接）
func isRegularFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.Mode().IsRegular()
}

// CheckServerFiles 检查多台服务器的服务端文件目录，返回各自的安装状态。
// 目录不存在的服务器不会出现在结果里
func (dm *DockerManager) CheckServerFiles(serverIDs []uint) (map[uint]ServerFilesState, error) {
	result := map[uint]ServerFilesState{}
	for _, id := range serverIDs {
		dir := ServerFilesDir(id)
		info, err := os.Stat(dir)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("检查服务端文件失败: %w", err)
		}
		if !info.IsDir() {
			continue
		}
		result[id] = ServerFilesState{
			Installed: isRegularFile(filepath.Join(dir, installMarkerFile)),
			Windows:   isRegularFile(filepath.Join(dir, filepath.FromSlash(windowsServerExecutable))),
			Linux:     isRegularFile(filepath.Join(dir, filepath.FromSlash(linuxServerExecutable))),
			Branch:    readGameBranch(dir),
		}
	}
	return result, nil
}

// CopyServerFiles 把 srcID 的服务端文件（含已安装的 Mod 和 .installed_mods）复制到 dstID 的文件卷，
// 并删除 .install_complete，让入口脚本启动时执行 steamcmd validate：只校验本地文件、下载有差异的部分
// （来源是另一个平台时也一样，validate 会补上这个平台的主程序）。ArkApi 文件在文件卷里，会一起复制。
// 存档和插件在单独的卷里，不会被复制
func (dm *DockerManager) CopyServerFiles(srcID, dstID uint) error {
	if srcID == dstID {
		return errors.New("来源和目标是同一台服务器")
	}
	if err := copyTree(ServerFilesDir(srcID), ServerFilesDir(dstID)); err != nil {
		return fmt.Errorf("复制服务端文件失败: %w", err)
	}
	if err := os.Remove(filepath.Join(ServerFilesDir(dstID), installMarkerFile)); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("复制服务端文件失败: %w", err)
	}
	return nil
}

// readGameBranch 读取文件卷记录的游戏版本，没有记录时为正式版
func readGameBranch(dir string) string {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return utils.GameBranchPublic
	}
	defer root.Close()
	data, err := root.ReadFile(gameBranchFile)
	if err != nil {
		return utils.GameBranchPublic
	}
	return strings.TrimSpace(string(data))
}

// PrepareGameBranch 启动前让文件卷和服务器的游戏版本一致：文件卷里已有服务端文件且属于另一个版本时，
// 删除 .install_complete（入口脚本会用新的 APP_UPDATE 重新 validate，只下载两个版本的差异），
// 切回正式版时还要去掉 appmanifest 里记录的 BetaKey，否则 steamcmd 会继续留在 beta 分支；
// 最后写入 .game_branch。文件卷目录不存在时不处理
func (dm *DockerManager) PrepareGameBranch(serverID uint, branch string) error {
	dir := ServerFilesDir(serverID)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("切换游戏版本失败: %w", err)
	}
	defer root.Close()

	if readGameBranch(dir) == branch {
		// 正式版的旧文件卷没有记录，补写一份
		if _, err := root.Stat(gameBranchFile); err == nil {
			return nil
		}
		return writeGameBranch(root, branch)
	}

	if err := root.Remove(installMarkerFile); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("切换游戏版本失败: %w", err)
	}
	if branch == utils.GameBranchPublic {
		if err := clearManifestBetaKey(root); err != nil {
			return fmt.Errorf("切换游戏版本失败: %w", err)
		}
	}
	return writeGameBranch(root, branch)
}

// clearManifestBetaKey 去掉 appmanifest 里的 BetaKey，steamcmd 下次 app_update 时回到默认分支。
// 没有 appmanifest（还没下载过）时不处理
func clearManifestBetaKey(root *os.Root) error {
	data, err := root.ReadFile(appManifestFile)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	cleaned := betaKeyLine.ReplaceAll(data, nil)
	if len(cleaned) == len(data) {
		return nil
	}
	return root.WriteFile(appManifestFile, cleaned, 0o644)
}

func writeGameBranch(root *os.Root, branch string) error {
	if err := root.WriteFile(gameBranchFile, []byte(branch+"\n"), 0o644); err != nil {
		return fmt.Errorf("记录游戏版本失败: %w", err)
	}
	return nil
}
