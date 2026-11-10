package utils

import (
	"slices"
	"strings"
)

// SetIniValue 在 INI 文本的指定 section 里设置一个键，返回新的文本。
// 键名和 section 名不区分大小写；已有的同名键（包括重复的）合并为一行，其它内容原样保留；
// section 不存在时追加到文件末尾。值里不能包含换行，由调用方保证
func SetIniValue(content, section, key, value string) string {
	var lines []string
	if strings.TrimSpace(content) != "" {
		lines = strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	}

	header := "[" + section + "]"
	entry := key + "=" + value
	out := make([]string, 0, len(lines)+3)
	inSection, found, replaced := false, false, false
	insertAt := -1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inSection = strings.EqualFold(trimmed, header)
			out = append(out, line)
			if inSection && !found {
				found = true
				insertAt = len(out)
			}
			continue
		}

		if inSection {
			if k, _, ok := strings.Cut(trimmed, "="); ok && strings.EqualFold(strings.TrimSpace(k), key) {
				if !replaced {
					out = append(out, entry)
					replaced = true
				}
				continue
			}
		}
		out = append(out, line)
	}

	if !replaced {
		if found {
			out = slices.Insert(out, insertAt, entry)
		} else {
			if len(out) > 0 && strings.TrimSpace(out[len(out)-1]) != "" {
				out = append(out, "")
			}
			out = append(out, header, entry)
		}
	}

	return strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
}
