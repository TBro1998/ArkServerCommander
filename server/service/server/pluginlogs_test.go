package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ark-server-commander/config"
)

func TestPluginLogs(t *testing.T) {
	setupTestDB(t)
	old := config.DataDir
	config.DataDir = t.TempDir()
	t.Cleanup(func() { config.DataDir = old })

	srv := createTestServer(t, 7777, 27015, 32330)
	id := fmt.Sprint(srv.ID)
	s := NewServerService()

	// 目录还不存在
	list, err := s.ListPluginLogs(1, id)
	if err != nil || list.Total != 0 || len(list.Files) != 0 {
		t.Fatalf("empty dir: %+v %v", list, err)
	}
	if _, err := s.GetPluginLog(1, id, "ArkApi.log", 10); !errors.Is(err, ErrPluginLogNotFound) {
		t.Fatalf("missing dir: %v", err)
	}

	dir := pluginLogsDir(srv.ID)
	if err := os.MkdirAll(filepath.Join(dir, "sub.log"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string, mod time.Time) {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(p, mod, mod); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	write("ArkApi_1_2026-11-23.log", "old\n", now.Add(-time.Hour))
	write("ArkApi_2_2026-11-24.log", "a\nb\nc\n", now)
	write("notes.txt", "x", now)
	write(".hidden.log", "x", now)
	outside := filepath.Join(config.DataDir, "secret.log")
	os.WriteFile(outside, []byte("secret"), 0o644)
	os.Symlink(outside, filepath.Join(dir, "link.log"))

	list, err = s.ListPluginLogs(1, id)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range list.Files {
		names = append(names, f.Name)
	}
	if got := strings.Join(names, ","); got != "ArkApi_2_2026-11-24.log,ArkApi_1_2026-11-23.log" || list.Total != 2 {
		t.Fatalf("files = %s (total %d)", got, list.Total)
	}

	log, err := s.GetPluginLog(1, id, "ArkApi_2_2026-11-24.log", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !log.Exists || strings.Join(log.Lines, "|") != "b|c" {
		t.Fatalf("log = %+v", log)
	}

	for _, name := range []string{"../secret.log", "notes.txt", ".hidden.log", "sub.log", "link.log", "nope.log", ""} {
		if _, err := s.GetPluginLog(1, id, name, 10); err == nil {
			t.Fatalf("%q should fail", name)
		}
	}
	if _, err := s.GetPluginLog(2, id, "ArkApi_2_2026-11-24.log", 10); err == nil {
		t.Fatal("other user should not read the log")
	}
}
