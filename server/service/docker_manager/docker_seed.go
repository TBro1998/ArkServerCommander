package docker_manager

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"ark-server-commander/utils"
)

// serverExecutable 服务端主程序相对安装目录的路径，存在才说明游戏本体已下载
const serverExecutable = "ShooterGame/Binaries/Win64/ShooterGameServer.exe"

// seedCopyTimeout 复制服务端文件（约 20 GB）的超时时间，机械硬盘上也留足余量
const seedCopyTimeout = 3 * time.Hour

// ServerFilesState 服务端文件卷的安装状态
type ServerFilesState struct {
	Installed  bool // 有 .install_complete 标记（入口脚本下载完成后写入）
	Executable bool // 服务端主程序存在
}

// checkFilesScript 逐个检查挂载在 /v/<id> 下的文件卷，每行输出 "<id> <有标记> <有主程序>"
const checkFilesScript = `for d in /v/*; do
  [ -d "$d" ] || continue
  m=0; e=0
  [ -f "$d/.install_complete" ] && m=1
  [ -f "$d/` + serverExecutable + `" ] && e=1
  echo "${d##*/} $m $e"
done`

// parseFilesStates 解析 checkFilesScript 的输出
func parseFilesStates(out string) map[uint]ServerFilesState {
	result := map[uint]ServerFilesState{}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		id, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			continue
		}
		result[uint(id)] = ServerFilesState{Installed: fields[1] == "1", Executable: fields[2] == "1"}
	}
	return result
}

// CheckServerFiles 用一个临时容器只读挂载多台服务器的服务端文件卷，返回各自的安装状态。
// 不存在的卷不会挂载（挂载会让 Docker 自动创建一个空卷），也不会出现在结果里
func (dm *DockerManager) CheckServerFiles(serverIDs []uint) (map[uint]ServerFilesState, error) {
	var binds []string
	for _, id := range serverIDs {
		name := utils.GetServerFilesVolumeName(id)
		exists, err := dm.VolumeExists(name)
		if err != nil {
			return nil, err
		}
		if exists {
			binds = append(binds, fmt.Sprintf("%s:/v/%d:ro", name, id))
		}
	}
	if len(binds) == 0 {
		return map[uint]ServerFilesState{}, nil
	}
	out, err := dm.RunHelper(binds, []string{"sh", "-c", checkFilesScript}, 2*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("检查服务端文件失败: %w", err)
	}
	return parseFilesStates(out), nil
}

// CopyServerFiles 把 srcID 的服务端文件（含已安装的 Mod 和 .installed_mods）复制到 dstID 的文件卷，
// 并删除 .install_complete，让入口脚本启动时执行 steamcmd validate：只校验本地文件、下载有差异的部分。
// 存档和插件在单独的卷里，不会被复制
func (dm *DockerManager) CopyServerFiles(srcID, dstID uint) error {
	binds := []string{
		fmt.Sprintf("%s:/src:ro", utils.GetServerFilesVolumeName(srcID)),
		fmt.Sprintf("%s:/dst", utils.GetServerFilesVolumeName(dstID)),
	}
	if _, err := dm.RunHelper(binds, []string{"sh", "-c", "cp -a /src/. /dst/ && rm -f /dst/.install_complete"}, seedCopyTimeout); err != nil {
		return fmt.Errorf("复制服务端文件失败: %w", err)
	}
	return nil
}
