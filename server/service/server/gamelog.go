package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf16"

	"ark-server-commander/service/docker_manager"
)

// gameLogReadLimit 读取游戏日志末尾的最大字节数
const gameLogReadLimit = 2 << 20

// GameLog 游戏自身的日志（存档目录 Logs/ShooterGame.log）
type GameLog struct {
	Exists     bool     `json:"exists"`
	Lines      []string `json:"lines"`
	Size       int64    `json:"size"`
	ModifiedAt string   `json:"modified_at,omitempty"`
}

// GetGameLog 返回 ShooterGame.log 最后 tail 行。直接读数据目录，不需要容器在运行
func (s *ServerService) GetGameLog(userID uint, serverID string, tail int) (*GameLog, error) {
	server, err := getUserServer(userID, serverID)
	if err != nil {
		return nil, err
	}
	p := filepath.Join(docker_manager.ServerSavedDir(server.ID), "Logs", "ShooterGame.log")
	f, err := os.Open(p)
	if errors.Is(err, os.ErrNotExist) {
		return &GameLog{Lines: []string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取游戏日志失败: %w", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("读取游戏日志失败: %w", err)
	}

	lines, err := tailLogFile(f, fi.Size(), tail)
	if err != nil {
		return nil, fmt.Errorf("读取游戏日志失败: %w", err)
	}
	return &GameLog{Exists: true, Lines: lines, Size: fi.Size(), ModifiedAt: fi.ModTime().Format(time.RFC3339)}, nil
}

// tailLogFile 读取日志末尾最多 gameLogReadLimit 字节，按编码解码后返回最后 tail 行。
// UE 的日志可能是 UTF-8，也可能是 UTF-16LE（带或不带 BOM），按 BOM 和零字节分布判断
func tailLogFile(r io.ReaderAt, size int64, tail int) ([]string, error) {
	head := make([]byte, 4)
	n, _ := r.ReadAt(head, 0)
	head = head[:n]
	utf16LE := bytes.HasPrefix(head, []byte{0xFF, 0xFE})

	offset := int64(0)
	if size > gameLogReadLimit {
		offset = size - gameLogReadLimit
	}
	buf := make([]byte, size-offset)
	if _, err := r.ReadAt(buf, offset); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if offset == 0 {
		buf = bytes.TrimPrefix(buf, []byte{0xEF, 0xBB, 0xBF})
		buf = bytes.TrimPrefix(buf, []byte{0xFF, 0xFE})
	}
	if !utf16LE && looksUTF16LE(buf, offset) {
		utf16LE = true
	}

	var text string
	if utf16LE {
		// UTF-16 每个码元 2 字节，从奇数偏移开始时丢掉第一个字节
		if offset%2 == 1 && len(buf) > 0 {
			buf = buf[1:]
		}
		u := make([]uint16, len(buf)/2)
		for i := range u {
			u[i] = uint16(buf[2*i]) | uint16(buf[2*i+1])<<8
		}
		text = string(utf16.Decode(u))
	} else {
		text = strings.ToValidUTF8(string(buf), "�")
	}

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	// 从中间开始读时第一行不完整
	if offset > 0 && len(lines) > 0 {
		lines = lines[1:]
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if tail > 0 && len(lines) > tail {
		lines = lines[len(lines)-tail:]
	}
	if lines == nil {
		lines = []string{}
	}
	return lines, nil
}

// looksUTF16LE 没有 BOM 时的判断：ASCII 为主的 UTF-16LE 文本，奇数位置几乎都是 0
func looksUTF16LE(buf []byte, offset int64) bool {
	start := int(offset % 2) // 让 i 对应文件中的偶数位置
	sample := buf
	if len(sample) > 4096 {
		sample = sample[:4096]
	}
	if len(sample) < 16 {
		return false
	}
	zeros, total := 0, 0
	for i := start + 1; i < len(sample); i += 2 {
		total++
		if sample[i] == 0 {
			zeros++
		}
	}
	return total > 0 && zeros*10 >= total*8
}
