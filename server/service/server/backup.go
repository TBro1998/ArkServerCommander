package server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"

	"ark-server-commander/config"
	"ark-server-commander/database"
	"ark-server-commander/models"
	"ark-server-commander/service/docker_manager"
	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// 备份类型，写在文件名里
const (
	BackupManual     = "manual"      // 手动备份
	BackupScheduled  = "scheduled"   // 定时任务
	BackupPreRestore = "pre-restore" // 恢复前自动备份当前存档
	BackupUploaded   = "uploaded"    // 用户上传
)

// backupNamePattern 备份文件名：<时间>_<类型>.tar.gz，例如 20261111-153000_manual.tar.gz
var backupNamePattern = regexp.MustCompile(`^(\d{8}-\d{6})_(manual|scheduled|pre-restore|uploaded)(?:-(\d+))?\.tar\.gz$`)

// BackupInfo 备份文件信息
type BackupInfo struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

// backupLocks 同一台服务器同时只能进行一个备份或恢复操作
var backupLocks sync.Map

func lockBackup(serverID uint) (func(), error) {
	mu, _ := backupLocks.LoadOrStore(serverID, &sync.Mutex{})
	m := mu.(*sync.Mutex)
	if !m.TryLock() {
		return nil, ErrBackupBusy
	}
	return m.Unlock, nil
}

// serverBackupDir 服务器的备份目录
func serverBackupDir(serverID uint) string {
	return filepath.Join(config.BackupDir, strconv.FormatUint(uint64(serverID), 10))
}

// newBackupPath 生成新备份的路径；同一秒内有重名时加序号
func newBackupPath(serverID uint, kind string) (string, string, error) {
	dir := serverBackupDir(serverID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", fmt.Errorf("创建备份目录失败: %w", err)
	}
	stamp := nowLocal().Format("20060102-150405")
	name := fmt.Sprintf("%s_%s.tar.gz", stamp, kind)
	for i := 2; ; i++ {
		if _, err := os.Stat(filepath.Join(dir, name)); errors.Is(err, os.ErrNotExist) {
			break
		}
		name = fmt.Sprintf("%s_%s-%d.tar.gz", stamp, kind, i)
	}
	return filepath.Join(dir, name), name, nil
}

// writeBackupFile 先写临时文件再改名，失败时不会留下不完整的备份
func writeBackupFile(path string, write func(io.Writer) error) error {
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %w", err)
	}
	if err := write(f); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("写入备份文件失败: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("保存备份文件失败: %w", err)
	}
	return nil
}

func backupInfo(dir, name string) (BackupInfo, bool) {
	m := backupNamePattern.FindStringSubmatch(name)
	if m == nil {
		return BackupInfo{}, false
	}
	st, err := os.Stat(filepath.Join(dir, name))
	if err != nil {
		return BackupInfo{}, false
	}
	created := st.ModTime()
	if t, err := time.ParseInLocation("20060102-150405", m[1], appLocation()); err == nil {
		created = t
	}
	return BackupInfo{Name: name, Kind: m[2], Size: st.Size(), CreatedAt: created.Format(time.RFC3339)}, true
}

// createBackup 备份服务器存档。运行中的服务器先通过 RCON 保存世界，保存失败也继续备份
func createBackup(server models.Server, kind string) (*BackupInfo, error) {
	unlock, err := lockBackup(server.ID)
	if err != nil {
		return nil, err
	}
	defer unlock()

	dockerManager, err := docker_manager.GetDockerManager()
	if err != nil {
		return nil, fmt.Errorf("获取Docker管理器失败: %w", err)
	}

	if server.Status == "running" {
		containerName := utils.GetServerContainerName(server.ID)
		if _, err := dockerManager.RconCommand(containerName, server.RCONPort, server.AdminPassword, "SaveWorld", time.Minute); err != nil {
			utils.Warn("备份前保存世界失败，继续备份", zap.Uint("server_id", server.ID), zap.Error(err))
		}
	}

	path, name, err := newBackupPath(server.ID, kind)
	if err != nil {
		return nil, err
	}
	if err := writeBackupFile(path, func(w io.Writer) error {
		return dockerManager.ExportSaved(server.ID, w)
	}); err != nil {
		return nil, err
	}

	info, _ := backupInfo(filepath.Dir(path), name)
	utils.Info("存档备份完成", zap.Uint("server_id", server.ID), zap.String("file", name), zap.Int64("size", info.Size))
	return &info, nil
}

// CreateBackup 手动备份
func (s *ServerService) CreateBackup(userID uint, serverID string) (*BackupInfo, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	if server.Status == "starting" || server.Status == "stopping" {
		return nil, ErrServerBusy
	}
	return createBackup(server, BackupManual)
}

// ListBackups 列出服务器的备份，最新的在前
func (s *ServerService) ListBackups(userID uint, serverID string) ([]BackupInfo, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	return listBackups(server.ID)
}

func listBackups(serverID uint) ([]BackupInfo, error) {
	dir := serverBackupDir(serverID)
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return []BackupInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %w", err)
	}
	result := []BackupInfo{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if info, ok := backupInfo(dir, e.Name()); ok {
			result = append(result, info)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name > result[j].Name })
	return result, nil
}

// BackupPath 返回备份文件的路径，用于下载。文件名必须符合备份命名规则，防止路径穿越
func (s *ServerService) BackupPath(userID uint, serverID, name string) (string, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return "", err
	}
	return existingBackupPath(server.ID, name)
}

func existingBackupPath(serverID uint, name string) (string, error) {
	if !backupNamePattern.MatchString(name) {
		return "", ErrBackupNotFound
	}
	path := filepath.Join(serverBackupDir(serverID), name)
	if st, err := os.Stat(path); err != nil || st.IsDir() {
		return "", ErrBackupNotFound
	}
	return path, nil
}

// DeleteBackup 删除一个备份
func (s *ServerService) DeleteBackup(userID uint, serverID, name string) error {
	path, err := s.BackupPath(userID, serverID, name)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("删除备份失败: %w", err)
	}
	return nil
}

// UploadBackup 保存用户上传的备份（本系统导出的 tar.gz），校验格式后放进备份列表
func (s *ServerService) UploadBackup(userID uint, serverID string, r io.Reader) (*BackupInfo, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	path, name, err := newBackupPath(server.ID, BackupUploaded)
	if err != nil {
		return nil, err
	}
	if err := writeBackupFile(path, func(w io.Writer) error {
		_, err := io.Copy(w, r)
		return err
	}); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("读取上传的备份失败: %w", err)
	}
	hasSaved, verr := docker_manager.ValidateBackup(f)
	f.Close()
	if verr == nil && !hasSaved {
		verr = fmt.Errorf("%w: 没有找到存档目录 Saved/SavedArks", ErrInvalidBackup)
	}
	if verr != nil {
		os.Remove(path)
		return nil, verr
	}

	info, _ := backupInfo(filepath.Dir(path), name)
	return &info, nil
}

// RestoreBackup 用备份覆盖服务器的存档和配置文件。只能在服务器停止时执行；
// 恢复前会自动备份当前存档（pre-restore），恢复后清除暂存的配置修改
func (s *ServerService) RestoreBackup(userID uint, serverID, name string) error {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return err
	}
	path, err := existingBackupPath(server.ID, name)
	if err != nil {
		return err
	}
	if server.Status != "stopped" {
		return ErrServerNotStopped
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("读取备份失败: %w", err)
	}
	_, verr := docker_manager.ValidateBackup(f)
	f.Close()
	if verr != nil {
		return verr
	}

	// 先做恢复前备份（内部会加锁），再占用服务器状态，防止恢复期间被启动
	if _, err := createBackup(server, BackupPreRestore); err != nil {
		return fmt.Errorf("恢复前备份当前存档失败: %w", err)
	}

	ok, err := transitionStatus(server.ID, []string{"running", "starting", "stopping"}, "stopping")
	if err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}
	if !ok {
		return ErrServerBusy
	}

	unlock, err := lockBackup(server.ID)
	if err != nil {
		finishOperation(server.ID, "stopped", "")
		return err
	}
	defer unlock()

	restoreErr := func() error {
		dockerManager, err := docker_manager.GetDockerManager()
		if err != nil {
			return fmt.Errorf("获取Docker管理器失败: %w", err)
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("读取备份失败: %w", err)
		}
		defer f.Close()
		return dockerManager.ImportSaved(server.ID, f)
	}()

	if restoreErr != nil {
		finishOperation(server.ID, "stopped", "恢复存档失败: "+restoreErr.Error())
		return restoreErr
	}

	database.DB.Model(&models.Server{}).Where("id = ?", server.ID).
		Updates(map[string]any{"pending_game_user_settings": "", "pending_game_ini": ""})
	finishOperation(server.ID, "stopped", "")
	utils.Info("存档已恢复", zap.Uint("server_id", server.ID), zap.String("file", name))
	return nil
}

// pruneBackups 只保留最近 keep 份指定类型的备份
func pruneBackups(serverID uint, kind string, keep int) {
	if keep <= 0 {
		return
	}
	backups, err := listBackups(serverID)
	if err != nil {
		return
	}
	count := 0
	for _, b := range backups {
		if b.Kind != kind {
			continue
		}
		count++
		if count > keep {
			if err := os.Remove(filepath.Join(serverBackupDir(serverID), b.Name)); err != nil {
				utils.Warn("清理旧备份失败", zap.String("file", b.Name), zap.Error(err))
			}
		}
	}
}

// backupDirSize 服务器备份目录占用的空间
func backupDirSize(serverID uint) int64 {
	backups, err := listBackups(serverID)
	if err != nil {
		return 0
	}
	var total int64
	for _, b := range backups {
		total += b.Size
	}
	return total
}
