package filewatcher

import (
	"fmt"
	"os"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	opts := []Option{
		OptID(""),
		OptFormat("json"),
		OptKind(Resource),
		OptName("Reader"),
		OptWatcher("", "json", "1s", false, false, "", true),
	}
	fw, err := NewFileWatcher(opts...) // already set global 'Event' to file watcher Event
	// fw.Event = &Event{}
	if err != nil {
		panic(err)
	}
	prepare := func(watcher *Watcher) {}
	cleanup := func(watcher *Watcher) {
		if err := os.RemoveAll(watcher.folder); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", watcher.folder)
	}
	fw.StartWait(prepare, cleanup)
}
