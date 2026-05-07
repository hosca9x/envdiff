package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestNew_LoadsInitialState(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	w, err := watcher.New(path, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Stop()
}

func TestNew_MissingFile(t *testing.T) {
	_, err := watcher.New("/nonexistent/.env", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempEnv(t, "FOO=original\n")
	w, err := watcher.New(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ch := w.Watch()

	// Give the watcher a moment before mutating the file.
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=changed\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Err != nil {
			t.Fatalf("event error: %v", ev.Err)
		}
		if len(ev.Entries) == 0 {
			t.Fatal("expected at least one diff entry")
		}
		if ev.File != path {
			t.Errorf("got file %q, want %q", ev.File, path)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
	w.Stop()
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnv(t, "STABLE=yes\n")
	w, err := watcher.New(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ch := w.Watch()

	select {
	case ev := <-ch:
		if ev.Err == nil {
			t.Fatalf("unexpected event for unchanged file: %+v", ev.Entries)
		}
	case <-time.After(150 * time.Millisecond):
		// expected: no event
	}
	w.Stop()
}

func TestStop_ClosesChannel(t *testing.T) {
	path := writeTempEnv(t, "X=1\n")
	w, err := watcher.New(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ch := w.Watch()
	w.Stop()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for channel close")
	}
}
