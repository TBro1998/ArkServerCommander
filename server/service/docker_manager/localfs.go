package docker_manager

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// copyTree 把 src 目录的内容合并复制到 dst（相当于 cp -a src/. dst/）：保留权限和修改时间，
// 符号链接原样复制（不跟随），其他特殊文件跳过；dst 里已有的同名文件被覆盖。
// 两边都用 os.Root 访问，目录里由 ARK 容器写入的符号链接不能把读写引到目录外
func copyTree(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	srcRoot, err := os.OpenRoot(src)
	if err != nil {
		return err
	}
	defer srcRoot.Close()
	dstRoot, err := os.OpenRoot(dst)
	if err != nil {
		return err
	}
	defer dstRoot.Close()

	type dirTime struct {
		rel   string
		perm  fs.FileMode
		mtime time.Time
	}
	var dirs []dirTime

	err = fs.WalkDir(srcRoot.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := filepath.FromSlash(p)
		info, err := d.Info()
		if err != nil {
			return err
		}
		switch {
		case d.IsDir():
			if p != "." {
				if err := replaceIfNot(dstRoot, rel, fs.ModeDir); err != nil {
					return err
				}
				if err := dstRoot.MkdirAll(rel, 0o755); err != nil {
					return err
				}
			}
			dirs = append(dirs, dirTime{rel, info.Mode().Perm(), info.ModTime()})
		case d.Type()&fs.ModeSymlink != 0:
			link, err := srcRoot.Readlink(rel)
			if err != nil {
				return err
			}
			if err := dstRoot.RemoveAll(rel); err != nil {
				return err
			}
			if err := dstRoot.Symlink(link, rel); err != nil {
				return fmt.Errorf("创建符号链接 %s 失败: %w", p, err)
			}
		case d.Type().IsRegular():
			if err := replaceIfNot(dstRoot, rel, 0); err != nil {
				return err
			}
			if err := copyFileInRoots(srcRoot, dstRoot, rel, info); err != nil {
				return fmt.Errorf("复制 %s 失败: %w", p, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 目录的权限和修改时间在写入子项后再设置（只读目录先设权限就写不进去，写入子项会改修改时间），从深到浅
	for _, d := range slices.Backward(dirs) {
		_ = dstRoot.Chmod(d.rel, d.perm)
		_ = dstRoot.Chtimes(d.rel, d.mtime, d.mtime)
	}
	return nil
}

// replaceIfNot dst 里 rel 已存在但类型不同（如要写文件的位置是目录或符号链接）时先删掉它。
// want 为 fs.ModeDir 表示目录，0 表示普通文件
func replaceIfNot(root *os.Root, rel string, want fs.FileMode) error {
	info, err := root.Lstat(rel)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode().Type() == want {
		return nil
	}
	return root.RemoveAll(rel)
}

// copyFileInRoots 复制一个普通文件，保留权限和修改时间。*os.File 之间的 io.Copy 在 Linux 上走 copy_file_range
func copyFileInRoots(srcRoot, dstRoot *os.Root, rel string, info fs.FileInfo) error {
	in, err := srcRoot.Open(rel)
	if err != nil {
		return err
	}
	defer in.Close()

	perm := info.Mode().Perm()
	out, err := dstRoot.OpenFile(rel, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm|0o200)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	_ = dstRoot.Chmod(rel, perm)
	_ = dstRoot.Chtimes(rel, info.ModTime(), info.ModTime())
	return nil
}
