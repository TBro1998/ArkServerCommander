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

// IniEntry INI 文件里的一个键值对
type IniEntry struct {
	Section string // section 名，不含方括号；文件开头没有 section 的内容为空
	Key     string
	Value   string
}

// iniKeyMatches 判断一行是不是 keys 里某个键（不区分大小写），注释行和非键值行不算
func iniKeyMatches(trimmed string, keys []string) (key, value string, ok bool) {
	if trimmed == "" || strings.HasPrefix(trimmed, ";") || strings.HasPrefix(trimmed, "#") {
		return "", "", false
	}
	k, v, found := strings.Cut(trimmed, "=")
	if !found {
		return "", "", false
	}
	k = strings.TrimSpace(k)
	for _, want := range keys {
		if strings.EqualFold(k, want) {
			return k, strings.TrimSpace(v), true
		}
	}
	return "", "", false
}

// FindIniEntries 返回所有 section 里键名属于 keys 的键值对（不区分大小写），按出现顺序
func FindIniEntries(content string, keys ...string) []IniEntry {
	var entries []IniEntry
	section := ""
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = trimmed[1 : len(trimmed)-1]
			continue
		}
		if k, v, ok := iniKeyMatches(trimmed, keys); ok {
			entries = append(entries, IniEntry{Section: section, Key: k, Value: v})
		}
	}
	return entries
}

// RemoveIniKeys 删除所有 section 里键名属于 keys 的行（不区分大小写），其它内容原样保留
func RemoveIniKeys(content string, keys ...string) string {
	if strings.TrimSpace(content) == "" {
		return content
	}
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if _, _, ok := iniKeyMatches(strings.TrimSpace(line), keys); ok {
			continue
		}
		out = append(out, line)
	}
	result := strings.TrimRight(strings.Join(out, "\n"), "\n")
	if strings.TrimSpace(result) == "" {
		return ""
	}
	return result + "\n"
}

// KeepIniKeys 以 content 为基础，把 keys 这些键换成 original 里的值（按 section 对应，不区分大小写）：
// content 里已有的行就地替换，original 里没有的删掉，content 里没有的追加到对应 section。
// 用于把配置文件复制到另一台服务器时保留它自己的名称、端口和密码
func KeepIniKeys(content, original string, keys ...string) string {
	id := func(section, key string) string { return strings.ToLower(section) + "]" + strings.ToLower(key) }
	originals := FindIniEntries(original, keys...)
	want := make(map[string]IniEntry, len(originals))
	for _, e := range originals {
		if _, ok := want[id(e.Section, e.Key)]; !ok {
			want[id(e.Section, e.Key)] = e
		}
	}

	placed := map[string]bool{}
	var out []string
	if strings.TrimSpace(content) != "" {
		section := ""
		for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
				section = trimmed[1 : len(trimmed)-1]
			} else if k, _, ok := iniKeyMatches(trimmed, keys); ok {
				key := id(section, k)
				if e, found := want[key]; found && !placed[key] {
					out = append(out, e.Key+"="+e.Value)
					placed[key] = true
				}
				continue
			}
			out = append(out, line)
		}
	}

	result := strings.TrimRight(strings.Join(out, "\n"), "\n")
	if strings.TrimSpace(result) == "" {
		result = ""
	} else {
		result += "\n"
	}
	for _, e := range originals {
		if key := id(e.Section, e.Key); !placed[key] && e.Section != "" {
			result = SetIniValue(result, e.Section, e.Key, e.Value)
			placed[key] = true
		}
	}
	return result
}
