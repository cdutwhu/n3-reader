package n3reader

import (
	"fmt"
	"os"
	"testing"

	fw "github.com/cdutwhu/n3-reader/file-watcher"
)

func TestNewN3Reader(t *testing.T) {
	if n3r, err := NewNats4Reader(); err == nil {
		opts := []fw.Option{
			fw.OptID(""),
			fw.OptFormat("json"),
			fw.OptName(""),
			fw.OptWatcher("", "json", "100ms", false, false, ""),
		}
		if freader, err := fw.NewFileWatcher(opts...); err == nil {
			freader.Event = NewN3ReaderEvent(n3r)
			cleanup := func(folder string) {
				if err := os.RemoveAll(folder); err != nil {
					panic(err)
				}
				fmt.Printf("%s is removed\n", folder)
			}
			freader.StartWait(cleanup)
		}
	}
}
