// Package watcher provides file-based change detection for .env files.
//
// It polls a target file at a configurable interval, computes a SHA-256
// checksum to detect mutations, and emits [Event] values containing the
// incremental diff whenever the file content changes.
//
// Basic usage:
//
//	w, err := watcher.New(".env", 5*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for ev := range w.Watch() {
//		if ev.Err != nil {
//			log.Println("watch error:", ev.Err)
//			continue
//		}
//		for _, e := range ev.Entries {
//			fmt.Println(e.Kind, e.Key)
//		}
//	}
//
// Call [Watcher.Stop] to terminate the polling goroutine and close the
// event channel.
package watcher
