// Package watcher monitors .env files for changes and emits diff events.
package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/loader"
)

// Event holds the diff result triggered by a file change.
type Event struct {
	File    string
	Entries []differ.Entry
	Err     error
}

// Watcher polls a file for changes and sends diff events on a channel.
type Watcher struct {
	file     string
	interval time.Duration
	lastHash [32]byte
	lastEnv  map[string]string
	mu       sync.Mutex
	stop     chan struct{}
}

// New creates a Watcher for the given file with the specified poll interval.
func New(file string, interval time.Duration) (*Watcher, error) {
	env, err := loader.LoadFile(file)
	if err != nil {
		return nil, fmt.Errorf("watcher: initial load: %w", err)
	}
	h, err := hashFile(file)
	if err != nil {
		return nil, fmt.Errorf("watcher: initial hash: %w", err)
	}
	return &Watcher{
		file:     file,
		interval: interval,
		lastHash: h,
		lastEnv:  env,
		stop:     make(chan struct{}),
	}, nil
}

// Watch starts polling and returns a channel that receives change events.
// Call Stop to terminate the watcher.
func (w *Watcher) Watch() <-chan Event {
	ch := make(chan Event, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				w.poll(ch)
			}
		}
	}()
	return ch
}

// Stop terminates the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll(ch chan<- Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	h, err := hashFile(w.file)
	if err != nil {
		ch <- Event{File: w.file, Err: err}
		return
	}
	if h == w.lastHash {
		return
	}
	newEnv, err := loader.LoadFile(w.file)
	if err != nil {
		ch <- Event{File: w.file, Err: err}
		return
	}
	entries := differ.Diff(w.lastEnv, newEnv)
	w.lastHash = h
	w.lastEnv = newEnv
	ch <- Event{File: w.file, Entries: entries}
}

func hashFile(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [32]byte{}, err
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out, nil
}
