package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "603.json")

	// The second write replaces an existing file.
	for _, content := range []string{"first", "second"} {
		if err := WriteFileAtomic(path, []byte(content), 0o644); err != nil {
			t.Fatalf("WriteFileAtomic(%q): %v", content, err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read back: %v", err)
		}
		if string(got) != content {
			t.Errorf("content = %q, want %q", got, content)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("dir holds %v, want only the target file (no temp leftovers)", names)
	}
}

func TestWriteFileAtomicMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "603.json")
	if err := WriteFileAtomic(path, []byte("x"), 0o644); err == nil {
		t.Fatal("want an error when the directory does not exist")
	}
}
