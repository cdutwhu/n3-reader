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
		OptName("Reader"),
		OptWatcher("", "json", "1s", false, false, ""),
	}
	fw, err := NewFileWatcher(opts...)
	if err != nil {
		panic(err)
	}
	prepare := func(watcher *Watcher) {}
	cleanup := func(watcher *Watcher) {
		if err := os.RemoveAll(watcher.Folder); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", watcher.Folder)
	}
	fw.StartWait(prepare, cleanup)
}
