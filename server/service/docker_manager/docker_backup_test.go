package docker_manager

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"testing"
)

type entry struct {
	name string
	typ  byte
	body string
}

func makeArchive(t *testing.T, entries []entry) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		h := &tar.Header{Name: e.name, Typeflag: e.typ, Mode: 0o644, Size: int64(len(e.body))}
		if e.typ == tar.TypeSymlink {
			h.Linkname = "/etc/passwd"
			h.Size = 0
		}
		if e.typ == tar.TypeDir {
			h.Size = 0
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if e.typ == tar.TypeReg {
			tw.Write([]byte(e.body))
		}
	}
	tw.Close()
	gz.Close()
	return &buf
}

func TestValidateBackup(t *testing.T) {
	good := makeArchive(t, []entry{
		{"Saved/", tar.TypeDir, ""},
		{"Saved/SavedArks/", tar.TypeDir, ""},
		{"Saved/SavedArks/TheIsland.ark", tar.TypeReg, "data"},
		{"Saved/Config/WindowsServer/Game.ini", tar.TypeReg, "[x]"},
	})
	has, err := ValidateBackup(good)
	if err != nil || !has {
		t.Fatalf("expected valid backup with SavedArks, got %v %v", has, err)
	}

	noSaves := makeArchive(t, []entry{{"Saved/Config/WindowsServer/Game.ini", tar.TypeReg, "[x]"}})
	if has, err := ValidateBackup(noSaves); err != nil || has {
		t.Fatalf("expected valid archive without SavedArks, got %v %v", has, err)
	}

	bad := map[string][]entry{
		"traversal": {{"Saved/../../etc/passwd", tar.TypeReg, "x"}},
		"outside":   {{"Other/file", tar.TypeReg, "x"}},
		"absolute":  {{"/Saved/file", tar.TypeReg, "x"}},
		"symlink":   {{"Saved/link", tar.TypeSymlink, ""}},
		"empty":     {},
	}
	for name, entries := range bad {
		if _, err := ValidateBackup(makeArchive(t, entries)); !errors.Is(err, ErrInvalidBackup) {
			t.Fatalf("%s: expected ErrInvalidBackup, got %v", name, err)
		}
	}

	if _, err := ValidateBackup(bytes.NewBufferString("not gzip")); !errors.Is(err, ErrInvalidBackup) {
		t.Fatalf("expected ErrInvalidBackup for non-gzip, got %v", err)
	}
}

func TestBackupSkipped(t *testing.T) {
	if !backupSkipped("Saved/Logs/ShooterGame.log") || !backupSkipped("Saved/clusters/abc/file") {
		t.Fatal("logs and clusters should be skipped")
	}
	if backupSkipped("Saved/SavedArks/TheIsland.ark") || backupSkipped("Saved/Config/WindowsServer/Game.ini") {
		t.Fatal("saves and config should be kept")
	}
}

func TestCPUPercent(t *testing.T) {
	prev := &StatsSample{CPUTotal: 1000, SystemCPU: 10000, OnlineCPUs: 4}
	cur := &StatsSample{CPUTotal: 2000, SystemCPU: 20000, OnlineCPUs: 4}
	if got := CPUPercent(prev, cur); got != 40 {
		t.Fatalf("expected 40, got %v", got)
	}
	if got := CPUPercent(nil, cur); got != 0 {
		t.Fatalf("expected 0 without previous sample, got %v", got)
	}
}

func TestParseFilesStates(t *testing.T) {
	out := "1 1 1\n2 0 1\n3 0 0\ngarbage\n4 1\n"
	got := parseFilesStates(out)
	want := map[uint]ServerFilesState{
		1: {Installed: true, Executable: true},
		2: {Installed: false, Executable: true},
		3: {},
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for id, st := range want {
		if got[id] != st {
			t.Fatalf("server %d: got %+v, want %+v", id, got[id], st)
		}
	}
}
