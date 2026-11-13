package server

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf16"
)

func encodeUTF16LE(s string, bom bool) []byte {
	var buf bytes.Buffer
	if bom {
		buf.Write([]byte{0xFF, 0xFE})
	}
	for _, u := range utf16.Encode([]rune(s)) {
		buf.WriteByte(byte(u))
		buf.WriteByte(byte(u >> 8))
	}
	return buf.Bytes()
}

func TestTailLogFile(t *testing.T) {
	text := "Log file open\r\n[2026.11.12] LogInit: 中文 line\r\nServer: \"x\" has successfully started!\r\n"
	want := []string{"[2026.11.12] LogInit: 中文 line", "Server: \"x\" has successfully started!"}

	for name, data := range map[string][]byte{
		"utf8":         []byte(text),
		"utf8 bom":     append([]byte{0xEF, 0xBB, 0xBF}, text...),
		"utf16 bom":    encodeUTF16LE(text, true),
		"utf16 no bom": encodeUTF16LE(text, false),
	} {
		t.Run(name, func(t *testing.T) {
			got, err := tailLogFile(bytes.NewReader(data), int64(len(data)), 2)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Join(got, "|") != strings.Join(want, "|") {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

func TestTailLogFileLarge(t *testing.T) {
	// 超过读取上限时从中间开始读，丢掉不完整的第一行
	var b strings.Builder
	for b.Len() < gameLogReadLimit+1000 {
		b.WriteString("0123456789 abcdefghij line\n")
	}
	b.WriteString("last line\n")
	for name, data := range map[string][]byte{
		"utf8":  []byte(b.String()),
		"utf16": encodeUTF16LE(b.String(), true),
	} {
		t.Run(name, func(t *testing.T) {
			got, err := tailLogFile(bytes.NewReader(data), int64(len(data)), 3)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 3 || got[2] != "last line" || got[0] != "0123456789 abcdefghij line" {
				t.Fatalf("unexpected tail %q", got)
			}
		})
	}
}
