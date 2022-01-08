package n3reader

import (
	"fmt"
	"os"
	"testing"

	fw "github.com/cdutwhu/n3-reader/file-watcher"
	lk "github.com/digisan/logkit"
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
	lk.FailOnErr("%v", err)

	opts := []Option{
		OptNatsHost(""),
		OptNatsPort(0),
		OptStream("ABC"),
		OptStreamSubjects("ABC.*"),
		OptSubject("ABC.created"),
	}
	nr, err := NewNats4Reader(opts...)
	lk.FailOnErr("%v", err)

	fw.Event = NewReaderEvent(nr)
	fw.StartWait(prepare, cleanup)
}
