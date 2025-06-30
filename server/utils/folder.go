package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// WorksDirectory 工作目录名称
	WorksDirectory = "Works"
)

// GetServerFolderPath 获取服务器文件夹路径
// serverID: 服务器数据库ID
// 返回: 服务器文件夹的绝对路径
func GetServerFolderPath(serverID uint) string {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		// 如果获取失败，使用当前目录
		currentDir = "."
	}

	// 构建路径: 当前目录/Works/服务器ID
	folderPath := filepath.Join(currentDir, WorksDirectory, fmt.Sprintf("%d", serverID))
	return folderPath
}

// CreateServerFolder 创建服务器文件夹
// serverID: 服务器数据库ID
// 返回: 创建的文件夹路径和错误信息
func CreateServerFolder(serverID uint) (string, error) {
	folderPath := GetServerFolderPath(serverID)

	// 创建文件夹（包括父目录）
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return "", fmt.Errorf("创建服务器文件夹失败: %v", err)
	}

	return folderPath, nil
}

// FolderExists 检查文件夹是否存在
// folderPath: 文件夹路径
// 返回: 文件夹是否存在
func FolderExists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// RemoveServerFolder 删除服务器文件夹
// serverID: 服务器数据库ID
// 返回: 错误信息
func RemoveServerFolder(serverID uint) error {
	folderPath := GetServerFolderPath(serverID)

	// 检查文件夹是否存在
	if !FolderExists(folderPath) {
		return nil // 文件夹不存在，视为删除成功
	}

	// 删除文件夹及其内容
	err := os.RemoveAll(folderPath)
	if err != nil {
		return fmt.Errorf("删除服务器文件夹失败: %v", err)
	}

	return nil
}

// GetFolderSize 获取文件夹大小（字节）
// folderPath: 文件夹路径
// 返回: 文件夹大小和错误信息
func GetFolderSize(folderPath string) (int64, error) {
	var size int64

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// EnsureWorksDirectory 确保Works目录存在
// 返回: Works目录路径和错误信息
func EnsureWorksDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	worksPath := filepath.Join(currentDir, WorksDirectory)

	// 创建Works目录
	err = os.MkdirAll(worksPath, 0755)
	if err != nil {
		return "", fmt.Errorf("创建Works目录失败: %v", err)
	}

	return worksPath, nil
}
