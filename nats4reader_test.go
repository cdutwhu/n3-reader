package n3reader

import (
	"fmt"
	"os"
	"testing"

	fw "github.com/cdutwhu/n3-reader/file-watcher"
)

func TestNewN3Reader(t *testing.T) {
	prepare := func(w *fw.Watcher) {

	}
	cleanup := func(w *fw.Watcher) {
		if err := os.RemoveAll(w.Folder()); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", w.Folder())
	}

	optsFR := []fw.Option{
		fw.OptID(""),
		fw.OptFormat("json"),
		fw.OptName(""),
		fw.OptWatcher("", "json", "100ms", false, false, "", true),
	}
	fw, err := fw.NewFileWatcher(optsFR...)
	Check(err)

	opts := []Option{
		OptNatsHost(""),
		OptNatsPort(0),
		OptStream("ABC"),
		OptStreamSubjects("ABC.*"),
		OptSubject("ABC.created"),
	}
	nr, err := NewNats4Reader(opts...)
	Check(err)

	fw.Event = NewReaderEvent(nr)
	fw.StartWait(prepare, cleanup)
}
