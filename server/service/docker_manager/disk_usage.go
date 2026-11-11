package docker_manager

import (
	"ark-server-commander/config"
	"ark-server-commander/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

// dirSize 统计目录下普通文件的总大小（字节），读不到的文件会被跳过
func dirSize(dir string) int64 {
	var total int64
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type().IsRegular() {
			if info, infoErr := d.Info(); infoErr == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total
}

// LocalVolumeSizes 统计数据目录下各服务器卷和集群卷占用的磁盘空间（字节），key 是卷名。
// 卷绑定在数据目录里，Docker 的 DiskUsage 查不到绑定卷的大小，所以直接遍历目录。
// 需要遍历全部文件，比较慢，调用方应缓存结果
func LocalVolumeSizes() (map[string]int64, error) {
	result := map[string]int64{}

	entries, err := os.ReadDir(filepath.Join(config.DataDir, "servers"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取服务器数据目录失败: %w", err)
	}
	for _, e := range entries {
		id, parseErr := strconv.ParseUint(e.Name(), 10, 32)
		if parseErr != nil || !e.IsDir() {
			continue
		}
		for _, spec := range serverVolumeSpecs(uint(id)) {
			result[spec.Name] = dirSize(spec.localDir())
		}
	}

	clusterEntries, err := os.ReadDir(filepath.Join(config.DataDir, "clusters"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取集群数据目录失败: %w", err)
	}
	for _, e := range clusterEntries {
		if !e.IsDir() {
			continue
		}
		spec := clusterVolumeSpec(e.Name())
		result[utils.GetClusterVolumeName(e.Name())] = dirSize(spec.localDir())
	}
	return result, nil
}
